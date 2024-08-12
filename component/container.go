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
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
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
	AddPostProcessor(processor ObjectProcessor) error
	PostProcessorCount() int
}

type ObjectContainer struct {
	definitions DefinitionRegistry
	singletons  SingletonRegistry

	scopes   map[string]Scope
	muScopes sync.RWMutex

	processors       []ObjectProcessor
	typesOfProcessor map[string]struct{}
	postProcessorMu  sync.RWMutex

	running bool
	mu      sync.RWMutex
}

func NewObjectContainer() *ObjectContainer {
	return &ObjectContainer{
		definitions:      NewObjectDefinitionRegistry(),
		singletons:       NewSingletonObjectRegistry(),
		scopes:           make(map[string]Scope),
		processors:       make([]ObjectProcessor, 0),
		typesOfProcessor: map[string]struct{}{},
	}
}

func (c *ObjectContainer) Start(ctx context.Context) error {
	defer c.mu.Unlock()
	c.mu.Lock()

	if c.running {
		return errors.New("container is already started")
	}

	err := c.loadObjectProcessors(ctx)
	if err != nil {
		return err
	}

	err = c.loadSingletonObjects(ctx)
	if err != nil {
		return err
	}

	c.running = true
	return nil
}

func (c *ObjectContainer) Stop(ctx context.Context) error {
	defer c.mu.Unlock()
	c.mu.Lock()

	if !c.running {
		return errors.New("container is already stopped")
	}

	c.running = false
	return nil
}

func (c *ObjectContainer) IsRunning() bool {
	defer c.mu.Unlock()
	c.mu.Lock()

	return false
}
func (c *ObjectContainer) GetObject(ctx context.Context, filters ...filter.Filter) (any, error) {
	ctx = contextWithHolder(ctx)

	candidate, err := c.Singletons().Find(filters...)
	if err == nil {
		return candidate, nil
	}

	var definition *Definition
	definition, err = c.Definitions().Find(filters...)

	if err != nil {
		return nil, err
	}

	objectName := definition.Name()

	if definition.IsSingleton() {
		var object any
		object, err = c.Singletons().OrElseCreate(objectName, func(ctx context.Context) (any, error) {

			if log.IsDebugEnabled() {
				log.D(ctx, "Creating singleton object of type '{}' under package '{}'", rawName(definition.Type()), definition.Type().PackagePath())
			}

			return c.createObject(ctx, definition, nil)
		})

		return object, err
	} else if definition.IsPrototype() {
		prototypeHolder := holderFromContext(ctx)
		err = prototypeHolder.beforeCreation(objectName)

		if err != nil {
			return nil, err
		}

		defer prototypeHolder.afterCreation(objectName)
		return c.createObject(ctx, definition, nil)
	}

	if strings.TrimSpace(definition.Scope()) == "" {
		return nil, fmt.Errorf("no scope name for required type %s", definition.Type().Name())
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
		return c.createObject(ctx, definition, nil)
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
	return c.singletons.Contains(name)
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
	for scopeName := range c.scopes {
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
		if reflector.IsStruct(objectType) && reflector.IsStruct(requiredType) {
			return matchTypeName(objectType, requiredType)
		}

		return true
	} else if reflector.IsPointer(objectType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
		ptrType := reflector.ToPointer(objectType)

		if ptrType.Elem().CanConvert(requiredType) {
			return true
		}
	}

	return false
}

func (c *ObjectContainer) AddObjectProcessor(processor ObjectProcessor) error {
	if processor == nil {
		return errors.New("processor cannot be nil")
	}

	defer c.postProcessorMu.Unlock()
	c.postProcessorMu.Lock()

	typ := reflector.TypeOfAny(processor)
	typeName := fullTypeName(typ)

	if _, ok := c.typesOfProcessor[typeName]; ok {
		return fmt.Errorf("processor '%s' is already registered", typeName)
	}

	c.typesOfProcessor[typeName] = struct{}{}
	c.processors = append(c.processors, processor)
	return nil
}

func (c *ObjectContainer) ObjectProcessorCount() int {
	defer c.postProcessorMu.Unlock()
	c.postProcessorMu.Lock()
	return len(c.processors)
}

func (c *ObjectContainer) createObject(ctx context.Context, definition *Definition, args []any) (object any, err error) {
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

		if reflector.TypeOfAny(results[0]).ReflectValue().IsZero() {
			return nil, fmt.Errorf("constructor function '%s' returns nil", constructorFunc.Name())
		}

		object = results[0]
	} else if (argsCount == 0 && len(args) == 0) || (len(args) != 0 && argsCount == len(args)) {
		var results []any
		results, err = constructorFunc.Invoke(args...)

		if err != nil {
			return nil, err
		}

		if reflector.TypeOfAny(results[0]).ReflectValue().IsZero() {
			return nil, fmt.Errorf("constructor function '%s' returns nil", constructorFunc.Name())
		}

		object = results[0]
	} else {
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition %s", definition.Name())
	}

	return c.initialize(ctx, object)
}

func (c *ObjectContainer) resolveArguments(ctx context.Context, args []*ConstructorArgument) ([]any, error) {
	arguments := make([]any, 0)

	for _, arg := range args {

		if reflector.IsSlice(arg.Type()) {
			sliceType := reflector.ToSlice(arg.Type())
			val, err := sliceType.Instantiate()

			if err != nil {
				return nil, err
			}

			sliceType = reflector.ToSlice(reflector.ToPointer(reflector.TypeOfAny(val.Val())).Elem())

			var result any
			objectList := c.ListObjects(ctx, filter.ByType(sliceType.Elem()))
			result, err = sliceType.Append(objectList...)

			arguments = append(arguments, result)
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

func (c *ObjectContainer) initialize(ctx context.Context, object any) (any, error) {
	result, err := c.applyProcessorsBeforeInit(ctx, object)
	if err != nil {
		return nil, err
	}

	if initialization, ok := object.(Initialization); ok {
		err = initialization.DoInit(ctx)

		if err != nil {
			return nil, err
		}
	}

	result, err = c.applyProcessorsAfterInit(ctx, object)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ObjectContainer) applyProcessorsBeforeInit(ctx context.Context, object any) (any, error) {
	for _, processor := range c.processors {
		result, err := processor.ProcessBeforeInit(ctx, object)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, fmt.Errorf("'%s' returns nil object from ProcessBeforeInit", reflector.TypeOfAny(processor).Name())
		}

		object = result
	}

	return object, nil
}

func (c *ObjectContainer) applyProcessorsAfterInit(ctx context.Context, object any) (any, error) {
	for _, processor := range c.processors {
		result, err := processor.ProcessAfterInit(ctx, object)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, fmt.Errorf("'%s' returns nil object from ProcessAfterInit", reflector.TypeOfAny(processor).Name())
		}

		object = result
	}

	return object, nil
}

func (c *ObjectContainer) loadObjectProcessors(ctx context.Context) error {

	postProcessors := c.Definitions().List(filter.ByTypeOf[ObjectProcessor]())

	checker := newProcessorChecker(c, c.ObjectProcessorCount()+len(postProcessors)+1)
	_ = c.AddObjectProcessor(checker)

	for _, processorDefinition := range postProcessors {
		processor, err := c.GetObject(ctx, filter.ByName(processorDefinition.Name()))

		if err != nil {
			return err
		}

		err = c.AddObjectProcessor(processor.(ObjectProcessor))

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ObjectContainer) loadSingletonObjects(ctx context.Context) error {

	for _, definition := range c.Definitions().List() {
		if !definition.IsSingleton() {
			continue
		}

		_, err := c.GetObject(ctx, filter.ByName(definition.Name()))

		if err != nil {
			return err
		}
	}

	return nil
}
