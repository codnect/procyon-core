package component

import (
	"codnect.io/reflector"
	"context"
	"errors"
	"fmt"
	"github.com/codnect/procyoncore/component/filter"
	"slices"
	"strings"
	"sync"
)

type Container interface {
	Start() error
	Stop() error
	IsRunning() bool
	CreateObject(ctx context.Context, definition *Definition, args []any) (any, error)
	GetObject(ctx context.Context, filters ...filter.Filter) (any, error)
	ListObjects(ctx context.Context, filters ...filter.Filter) []any
	ContainsObject(name string) bool
	IsSingleton(name string) bool
	IsPrototype(name string) bool
	Definitions() DefinitionRegistry
	Singletons() SingletonRegistry
	RegisterScope(name string, scope Scope) error
	ScopeNames() []string
	FindScope(name string) (Scope, error)
}

type ObjectContainer struct {
	definitions DefinitionRegistry
	singletons  SingletonRegistry

	scopes   map[string]Scope
	muScopes *sync.RWMutex
}

func NewObjectContainer() *ObjectContainer {
	return &ObjectContainer{
		definitions: NewObjectDefinitionRegistry(),
		singletons:  NewSingletonObjectRegistry(),
		scopes:      make(map[string]Scope),
		muScopes:    &sync.RWMutex{},
	}
}

func (c *ObjectContainer) Start() error {
	return nil
}

func (c *ObjectContainer) Stop() error {
	return nil
}

func (c *ObjectContainer) IsRunning() bool {
	return false
}
func (c *ObjectContainer) GetObject(ctx context.Context, filters ...filter.Filter) (any, error) {
	filterOpts := filter.Of(filters...)

	ctx = contextWithHolder(ctx)
	objectName := filterOpts.Name

	if objectName == "" && filterOpts.Type == nil {
		return nil, errors.New("filtering should be done by either name or type")
	}

	if objectName == "" {
		candidate, err := c.Singletons().Find(filter.ByType(filterOpts.Type))

		if err == nil {
			return candidate, nil
		}

		definitionList := c.Definitions().List(filter.ByType(filterOpts.Type))

		if len(definitionList) == 0 {
			return nil, ObjectNotFoundError{
				typ: filterOpts.Type,
			}
		} else if len(definitionList) > 1 {
			return nil, fmt.Errorf("there is more than one definition for the required type %s, it cannot be distinguished", filterOpts.Type.Name())
		}

		objectName = definitionList[0].Name()
	}

	definition, err := c.Definitions().Find(filter.ByName(objectName))

	if err != nil {
		return nil, err
	}

	if filterOpts.Type != nil && !c.matchType(definition.Type(), filterOpts.Type) {
		return nil, fmt.Errorf("definition type with name %s does not match the required type", objectName)
	}

	if definition.IsSingleton() {
		var object any
		object, err = c.Singletons().OrElseCreate(objectName, func(ctx context.Context) (any, error) {

			if log.IsDebugEnabled() {
				log.D(ctx, "Creating singleton object of '{}' under package '{}'", definition.Type().Name(), definition.Type().PackagePath())
			}

			return c.CreateObject(ctx, definition, nil)
		})

		return object, err
	} else if definition.IsPrototype() {
		prototypeHolder := holderFromContext(ctx)
		err = prototypeHolder.beforeCreation(objectName)

		if err != nil {
			return nil, err
		}

		defer prototypeHolder.afterCreation(objectName)
		return c.CreateObject(ctx, definition, nil)
	}

	if strings.TrimSpace(definition.Scope()) == "" {
		return nil, fmt.Errorf("no scope name for required type %s", filterOpts.Type.Name())
	}

	var scope Scope
	scope, err = c.FindScope(definition.Scope())

	if err != nil {
		return nil, err
	}

	return scope.GetObject(ctx, objectName, func(ctx context.Context) (any, error) {
		scopeHolder := holderFromContext(ctx)
		err = scopeHolder.beforeCreation(objectName)

		if err != nil {
			return nil, err
		}

		defer scopeHolder.afterCreation(objectName)
		return c.CreateObject(ctx, definition, nil)
	})
}

func (c *ObjectContainer) ListObjects(ctx context.Context, filters ...filter.Filter) []any {
	objectList := make([]any, 0)
	singletonNames := c.singletons.Names()
	objectList = append(objectList, c.singletons.List(filters...)...)

	for _, definition := range c.definitions.List(filters...) {
		if (definition.IsSingleton() && !slices.Contains(singletonNames, definition.Name())) || !definition.IsSingleton() {
			object, err := c.GetObject(ctx, filter.ByName(definition.Name()))

			if err != nil {
				continue
			}

			objectList = append(objectList, object)
		}
	}

	return objectList
}

func (c *ObjectContainer) ContainsObject(name string) bool {
	return true
}

func (c *ObjectContainer) IsSingleton(name string) bool {
	definition, ok := c.definitions.FindFirst(filter.ByName(name))
	return ok && definition.IsSingleton()
}

func (c *ObjectContainer) IsPrototype(name string) bool {
	definition, ok := c.definitions.FindFirst(filter.ByName(name))
	return ok && definition.IsPrototype()
}

func (c *ObjectContainer) Definitions() DefinitionRegistry {
	return c.definitions
}

func (c *ObjectContainer) Singletons() SingletonRegistry {
	return c.singletons
}

func (c *ObjectContainer) RegisterScope(name string, scope Scope) error {
	if strings.TrimSpace(name) == "" {
		panic("name cannot be empty or blank")
	}

	if scope == nil {
		panic("scope cannot be nil")
	}

	if SingletonScope != name && PrototypeScope != name {
		defer c.muScopes.Unlock()
		c.muScopes.Lock()
		c.scopes[name] = scope
		return nil
	}

	return errors.New("singleton' and 'prototype' cannot be replaced")

}

func (c *ObjectContainer) ScopeNames() []string {
	defer c.muScopes.Unlock()
	c.muScopes.Lock()

	scopeNames := make([]string, 0)
	for scopeName, _ := range c.scopes {
		scopeNames = append(scopeNames, scopeName)
	}

	return scopeNames
}

func (c *ObjectContainer) FindScope(name string) (Scope, error) {
	defer c.muScopes.Unlock()
	c.muScopes.Lock()
	if scope, ok := c.scopes[name]; ok {
		return scope, nil
	}

	return nil, fmt.Errorf("no scope registered for scope name %s", name)
}

func (c *ObjectContainer) matchType(objectType reflector.Type, requiredType reflector.Type) bool {
	if objectType.CanConvert(requiredType) {
		return true
	} else if reflector.IsPointer(objectType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
		ptrType := reflector.ToPointer(objectType)

		if ptrType.Elem().CanConvert(requiredType) {
			return true
		}
	}

	return false
}

func (c *ObjectContainer) CreateObject(ctx context.Context, definition *Definition, args []any) (object any, err error) {
	if ctx == nil {
		return nil, errors.New("ctx cannot be nil")
	}

	if definition == nil {
		return nil, errors.New("definition cannot be nil")
	}

	constructorFunc := definition.constructor
	argsCount := len(definition.ConstructorArguments())

	if argsCount != 0 && len(args) == 0 {
		var resolvedArguments []any
		resolvedArguments, err = c.resolveArguments(ctx, definition.ConstructorArguments())

		if err != nil {
			return nil, err
		}

		var results []any
		results, err = constructorFunc.Invoke(resolvedArguments...)

		if err != nil {
			return nil, err
		}

		object = results[0]
	} else if (argsCount == 0 && len(args) == 0) || (len(args) != 0 && argsCount == len(args)) {
		var results []any
		results, err = constructorFunc.Invoke(args...)

		if err != nil {
			return nil, err
		}

		object = results[0]
	} else {
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition %s", definition.Name())
	}

	//return c.initializeInstance(definition.name, instance)
	return
}

func (c *ObjectContainer) resolveArguments(ctx context.Context, args []*ConstructorArgument) ([]any, error) {
	arguments := make([]any, 0)

	for _, arg := range args {

		if reflector.IsSlice(arg.Type()) {
			sliceType := reflector.ToSlice(arg.Type())
			objectList := c.ListObjects(ctx, filter.ByType(sliceType.Elem()))
			arguments = append(arguments, objectList)
			continue
		}

		var (
			object any
			err    error
		)

		/*
			resolvableInstance, exists := m.getResolvableInstance(arg.Type())
			if exists {
				arguments = append(arguments, resolvableInstance)
				continue
			}
		*/

		if arg.Name() != "" {
			object, err = c.GetObject(ctx, filter.ByName(arg.Name()))
		} else {
			object, err = c.GetObject(ctx, filter.ByType(arg.Type()))
		}

		if err != nil {
			if notFoundErr := (*ObjectNotFoundError)(nil); errors.As(err, &notFoundErr) && !reflector.IsPointer(arg.Type()) && arg.Type().IsInstantiable() {
				var val reflector.Value
				val, err = arg.Type().Instantiate()

				if err != nil {
					return nil, err
				}

				object = val.Elem()
				arguments = append(arguments, object)
				continue
			}

			if !arg.IsOptional() && err != nil {
				return nil, err
			} else if arg.IsOptional() && err != nil {
				arguments = append(arguments, nil)
			}
		} else {
			arguments = append(arguments, object)
		}
	}

	return arguments, nil
}
