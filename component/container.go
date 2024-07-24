package component

import (
	"codnect.io/reflector"
	"context"
	"errors"
	"fmt"
	"github.com/codnect/procyoncore/component/filter"
	"strings"
	"sync"
)

type Container interface {
	Start() error
	Stop() error
	GetObject(ctx context.Context, filters ...filter.Filter) (any, error)
	ListObjects(ctx context.Context, filters ...filter.Filter) ([]any, error)
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

func (c *ObjectContainer) GetObject(ctx context.Context, filters ...filter.Filter) (any, error) {
	filterOpts := &filter.Filters{}
	for _, filterFunction := range filters {
		filterFunction(filterOpts)
	}

	ctx = contextWithHolder(ctx)
	objectName := filterOpts.Name

	if objectName == "" && filterOpts.Type == nil {
		return nil, errors.New("filtering should be done by either name or type")
	}

	if objectName == "" {
		candidate, err := c.Singletons().FindByType(filterOpts.Type)

		if err == nil {
			return candidate, nil
		}

		candidateNames := c.Definitions().NamesByType(filterOpts.Type)

		if len(candidateNames) == 0 {
			return nil, &NotFoundError{
				//ErrorString: fmt.Sprintf("container: not found instance or definition with required type %s", requiredType.Name()),
			}
		} else if len(candidateNames) > 1 {
			return nil, fmt.Errorf("there is more than one definition for the required type %s, it cannot be distinguished", filterOpts.Type.Name())
		}

		objectName = candidateNames[0]
	}

	definition, ok := c.Definitions().Find(objectName)

	if !ok {
		return nil, &NotFoundError{
			//ErrorString: fmt.Sprintf("container: not found definition with name %s", name),
		}
	}

	if filterOpts.Type != nil && !c.matchType(definition.Type(), filterOpts.Type) {
		return nil, fmt.Errorf("definition type with name %s does not match the required type", objectName)
	}

	if definition.IsSingleton() {
		instance, err := c.Singletons().OrElseGet(objectName, func(ctx context.Context) (any, error) {

			if log.IsDebugEnabled() {
				log.D(ctx, "Creating singleton instance of '{}' under package '{}'", definition.Type().Name(), definition.Type().PackagePath())
			}

			return c.createInstance(ctx, definition, filterOpts.Arguments)
		})

		return instance, err
	} else if definition.IsPrototype() {
		prototypeHolder := holderFromContext(ctx)
		err := prototypeHolder.beforeCreation(objectName)

		if err != nil {
			return nil, err
		}

		defer prototypeHolder.afterCreation(objectName)
		return c.createInstance(ctx, definition, filterOpts.Arguments)
	}

	if strings.TrimSpace(definition.Scope()) == "" {
		return nil, fmt.Errorf("no scope name for required type %s", filterOpts.Type.Name())
	}

	scope, err := c.FindScope(definition.Scope())

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
		return c.createInstance(ctx, definition, filterOpts.Arguments)
	})
}

func (c *ObjectContainer) ListObjects(ctx context.Context, filters ...filter.Filter) ([]any, error) {
	filterOpts := &filter.Filters{}
	for _, filterFunction := range filters {
		filterFunction(filterOpts)
	}

	return nil, nil
}

func (c *ObjectContainer) IsSingleton(name string) bool {
	def, ok := c.definitions.Find(name)
	return ok && def.IsSingleton()
}

func (c *ObjectContainer) IsPrototype(name string) bool {
	def, ok := c.definitions.Find(name)
	return ok && def.IsPrototype()
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

func (c *ObjectContainer) matchType(instanceType reflector.Type, requiredType reflector.Type) bool {
	if instanceType.CanConvert(requiredType) {
		return true
	} else if reflector.IsPointer(instanceType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
		ptrType := reflector.ToPointer(instanceType)

		if ptrType.Elem().CanConvert(requiredType) {
			return true
		}
	}

	return false
}

func (c *ObjectContainer) createInstance(ctx context.Context, definition *Definition, args []any) (instance any, err error) {
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

		instance = results[0]
	} else if (argsCount == 0 && len(args) == 0) || (len(args) != 0 && argsCount == len(args)) {
		var results []any
		results, err = constructorFunc.Invoke(args...)

		if err != nil {
			return nil, err
		}

		instance = results[0]
	} else {
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition %s", definition.Name())
	}

	//return c.initializeInstance(definition.name, instance)
	return
}

func (m *ObjectContainer) resolveArguments(ctx context.Context, args []*ConstructorArgument) ([]any, error) {
	arguments := make([]any, 0)

	for _, arg := range args {

		if reflector.IsSlice(arg.Type()) {
			/*sliceType := reflector.ToSlice(arg.Type())
			instances, err := m.getInstances(ctx, sliceType, sliceType.Elem())

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, instances)
			continue*/
		}

		var (
			instance any
			err      error
		)

		/*
			resolvableInstance, exists := m.getResolvableInstance(arg.Type())
			if exists {
				arguments = append(arguments, resolvableInstance)
				continue
			}
		*/

		if arg.Name() != "" {
			instance, err = m.GetObject(ctx, filter.ByName(arg.Name()))
		} else {
			instance, err = m.GetObject(ctx, filter.ByType(arg.Type()))
		}

		if err != nil {
			if notFoundErr := (*NotFoundError)(nil); errors.As(err, &notFoundErr) && !reflector.IsPointer(arg.Type()) && arg.Type().IsInstantiable() {
				var val reflector.Value
				val, err = arg.Type().Instantiate()

				if err != nil {
					return nil, err
				}

				instance = val.Elem()
				arguments = append(arguments, instance)
				continue
			}

			if !arg.IsOptional() && err != nil {
				return nil, err
			} else if arg.IsOptional() && err != nil {
				arguments = append(arguments, nil)
			}
		} else {
			arguments = append(arguments, instance)
		}
	}

	return arguments, nil
}
