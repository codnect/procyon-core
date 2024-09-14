package component

import (
	"codnect.io/procyon-core/component/filter"
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type Container interface {
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
	AddObjectProcessor(processor ObjectProcessor) error
	ObjectProcessorCount() int
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

func (c *ObjectContainer) GetObject(ctx context.Context, filters ...filter.Filter) (any, error) {
	if len(filters) == 0 {
		return nil, errors.New("at least one filter must be used")
	}

	ctx = contextWithHolder(ctx)

	candidate, err := c.Singletons().Find(filters...)
	if err == nil {
		return candidate, nil
	} else if _, ok := err.(*ObjectNotFoundError); !ok {
		return nil, err
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
				log.D(ctx, "Creating singleton object of type '{}'", definition.Type().String())
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
		panic("cannot register scope with empty or blank name")
	}

	if scope == nil {
		panic("nil scope")
	}

	if SingletonScope != name && PrototypeScope != name {
		defer c.muScopes.Unlock()
		c.muScopes.Lock()
		c.scopes[name] = scope
		return nil
	}

	return errors.New("cannot replace 'singleton' and 'prototype' scopes")

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

	return nil, fmt.Errorf("no scope registered for scope name '%s'", name)
}

func (c *ObjectContainer) AddObjectProcessor(processor ObjectProcessor) error {
	if processor == nil {
		return errors.New("nil processor")
	}

	defer c.postProcessorMu.Unlock()
	c.postProcessorMu.Lock()

	typ := reflect.TypeOf(processor)
	typeName := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())

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
		return nil, errors.New("nil context")
	}

	if definition == nil {
		return nil, errors.New("nil definition")
	}

	constructor := definition.Constructor()
	argsCount := len(constructor.Arguments())

	if argsCount != 0 && len(args) == 0 {
		var resolvedArguments []any
		resolvedArguments, err = c.resolveArguments(ctx, constructor.Arguments())

		if err != nil {
			return nil, err
		}

		var results []any
		results, err = constructor.Invoke(resolvedArguments...)

		if err != nil {
			return nil, err
		}

		resultType := reflect.TypeOf(results[0])
		resultValue := reflect.ValueOf(results[0])
		if (resultType.Kind() == reflect.Pointer || resultType.Kind() == reflect.Interface) && resultValue.IsZero() {
			return nil, fmt.Errorf("constructor function '%s' returns nil", constructor.Name())
		}

		object = results[0]
	} else if (argsCount == 0 && len(args) == 0) || (len(args) != 0 && argsCount == len(args)) {
		var results []any
		results, err = constructor.Invoke(args...)

		if err != nil {
			return nil, err
		}

		resultType := reflect.TypeOf(results[0])
		resultValue := reflect.ValueOf(results[0])

		if (resultType.Kind() == reflect.Pointer || resultType.Kind() == reflect.Interface) && resultValue.IsZero() {
			return nil, fmt.Errorf("constructor function '%s' returns nil", constructor.Name())
		}

		object = results[0]
	} else {
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition '%s'", definition.Name())
	}

	return c.initialize(ctx, object)
}

func (c *ObjectContainer) resolveArguments(ctx context.Context, args []ConstructorArgument) ([]any, error) {
	arguments := make([]any, 0)

	for _, arg := range args {

		if arg.Type().Kind() == reflect.Slice {
			sliceType := arg.Type()
			sliceVal := reflect.MakeSlice(sliceType, 0, 0)

			objectList := c.ListObjects(ctx, filter.ByType(sliceType.Elem()))
			for _, object := range objectList {
				sliceVal = reflect.Append(sliceVal, reflect.ValueOf(object))
			}

			arguments = append(arguments, sliceVal.Interface())
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
			argKind := arg.Type().Kind()

			if _, ok := err.(*ObjectNotFoundError); ok && argKind != reflect.Pointer && argKind != reflect.Interface {
				var val reflect.Value
				val = reflect.New(arg.Type())

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
			return nil, fmt.Errorf("'%s' returns nil object from ProcessBeforeInit", reflect.TypeOf(processor).Name())
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
			return nil, fmt.Errorf("'%s' returns nil object from ProcessAfterInit", reflect.TypeOf(processor).Name())
		}

		object = result
	}

	return object, nil
}

func (c *ObjectContainer) loadObjectProcessors(ctx context.Context) error {

	/*
		postProcessors := c.Definitions().List(filter.ByTypeOf[component.ObjectProcessor]())

		checker := component.newProcessorChecker(c, c.ObjectProcessorCount()+len(postProcessors)+1)
		_ = c.AddObjectProcessor(checker)

		for _, processorDefinition := range postProcessors {
			processor, err := c.GetObject(ctx, filter.ByName(processorDefinition.Name()))

			if err != nil {
				return err
			}

			err = c.AddObjectProcessor(processor.(component.ObjectProcessor))

			if err != nil {
				return err
			}
		}*/

	return nil
}
