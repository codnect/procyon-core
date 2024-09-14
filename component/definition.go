package component

import (
	"codnect.io/procyon-core/component/filter"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Option func(definition *Definition) error

type Definition struct {
	name        string
	typ         reflect.Type
	scope       string
	constructor *Constructor
}

type DefinitionRegistry interface {
	Register(definition *Definition) error
	Remove(name string) error
	Contains(name string) bool
	Find(filters ...filter.Filter) (*Definition, error)
	FindFirst(filters ...filter.Filter) (*Definition, bool)
	List(filters ...filter.Filter) []*Definition
	Names() []string
	Count() int
}

func MakeDefinition(constructorFunc ConstructorFunc, options ...Option) (*Definition, error) {
	if constructorFunc == nil {
		return nil, fmt.Errorf("nil constructor function")
	}

	constructorType := reflect.TypeOf(constructorFunc)

	if constructorType.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructor must be a function")
	}

	if constructorType.NumOut() != 1 {
		return nil, fmt.Errorf("constructor must only be a function returning one result")
	}

	returnType := constructorType.Out(0)
	definitionName := ""

	if returnType.Kind() == reflect.Pointer {
		definitionName = returnType.Elem().Name()
	} else {
		definitionName = returnType.Name()
	}

	definition := &Definition{
		name:  lowerCamelCase(definitionName),
		typ:   returnType,
		scope: SingletonScope,
		constructor: &Constructor{
			funcType:  constructorType,
			funcValue: reflect.ValueOf(constructorFunc),
			arguments: make([]ConstructorArgument, 0),
		},
	}

	numIn := constructorType.NumIn()
	constructor := definition.Constructor()

	for index := 0; index < numIn; index++ {
		argType := constructorType.In(index)

		arg := ConstructorArgument{
			index:    index,
			typ:      argType,
			optional: false,
		}

		constructor.arguments = append(constructor.arguments, arg)
	}

	for _, option := range options {
		err := option(definition)
		if err != nil {
			return nil, err
		}
	}

	return definition, nil
}

func (d *Definition) Name() string {
	return d.name
}

func (d *Definition) Type() reflect.Type {
	return d.typ
}

func (d *Definition) Constructor() *Constructor {
	return d.constructor
}

func (d *Definition) Scope() string {
	return d.scope
}

func (d *Definition) IsSingleton() bool {
	return d.scope == SingletonScope
}

func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}

type ObjectDefinitionRegistry struct {
	definitionMap map[string]*Definition
	muDefinitions *sync.RWMutex
}

func NewObjectDefinitionRegistry() *ObjectDefinitionRegistry {
	return &ObjectDefinitionRegistry{
		definitionMap: map[string]*Definition{},
		muDefinitions: &sync.RWMutex{},
	}
}

func (r *ObjectDefinitionRegistry) Register(definition *Definition) error {
	if definition == nil {
		return fmt.Errorf("nil definition")
	}

	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[definition.Name()]; exists {
		return fmt.Errorf("definition with name '%s' already exists", definition.Name())
	}

	r.definitionMap[definition.Name()] = definition

	return nil
}

func (r *ObjectDefinitionRegistry) Remove(name string) error {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[name]; !exists {
		return fmt.Errorf("no found definition with name '%s'", name)
	}

	delete(r.definitionMap, name)
	return nil
}

func (r *ObjectDefinitionRegistry) Contains(name string) bool {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	_, exists := r.definitionMap[name]
	return exists
}

func (r *ObjectDefinitionRegistry) Find(filters ...filter.Filter) (*Definition, error) {
	if len(filters) == 0 {
		return nil, errors.New("at least one filter must be used")
	}

	definitionList := r.List(filters...)

	if len(definitionList) > 1 {
		return nil, errors.New("cannot distinguish definitions because too many matching found")
	}

	if len(definitionList) == 0 {
		filterOpts := filter.Of(filters...)

		return nil, DefinitionNotFoundError{
			name: filterOpts.Name,
			typ:  filterOpts.Type,
		}
	}

	return definitionList[0], nil
}

func (r *ObjectDefinitionRegistry) FindFirst(filters ...filter.Filter) (*Definition, bool) {
	definitionList := r.List(filters...)

	if len(definitionList) == 0 {
		return nil, false
	}

	return definitionList[0], true
}

func (r *ObjectDefinitionRegistry) List(filters ...filter.Filter) []*Definition {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	filterOpts := filter.Of(filters...)
	definitionList := make([]*Definition, 0)

	for _, definition := range r.definitionMap {
		if filterOpts.Name != "" && filterOpts.Name != definition.Name() {
			continue
		}

		if filterOpts.Type == nil {
			definitionList = append(definitionList, definition)
			continue
		}

		if matchTypes(definition.Type(), filterOpts.Type) {
			definitionList = append(definitionList, definition)
		}
	}

	return definitionList
}

func (r *ObjectDefinitionRegistry) Names() []string {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	names := make([]string, 0)
	for name := range r.definitionMap {
		names = append(names, name)
	}

	return names
}

func (r *ObjectDefinitionRegistry) Count() int {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	return len(r.definitionMap)
}
