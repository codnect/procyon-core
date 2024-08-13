package component

import (
	"codnect.io/reflector"
	"errors"
	"fmt"
	"github.com/codnect/procyoncore/component/filter"
	"strings"
	"sync"
	"unicode"
)

type Option func(definition *Definition) error

type Definition struct {
	name            string
	typ             reflector.Type
	scope           string
	constructor     reflector.Function
	constructorArgs []*ConstructorArgument
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

func MakeDefinition(constructor Constructor, options ...Option) (*Definition, error) {
	if constructor == nil {
		return nil, fmt.Errorf("constructor should not be nil")
	}

	constructorType := reflector.TypeOfAny(constructor)
	if !reflector.IsFunction(constructorType) {
		return nil, fmt.Errorf("constructor should be a function")
	}

	constructorFunc := reflector.ToFunction(constructorType)
	if constructorFunc.NumResult() != 1 {
		return nil, fmt.Errorf("constructor can only be a function returning one result")
	}

	returnType := constructorFunc.Results()[0]
	definitionName := ""

	if reflector.IsPointer(returnType) {
		pointerType := reflector.ToPointer(returnType)
		definitionName = pointerType.Elem().Name()
	} else {
		definitionName = returnType.Name()
	}

	definition := &Definition{
		name:            lowerCamelCase(definitionName),
		typ:             returnType,
		scope:           SingletonScope,
		constructor:     constructorFunc,
		constructorArgs: make([]*ConstructorArgument, 0),
	}

	for index, parameterType := range constructorFunc.Parameters() {
		arg := &ConstructorArgument{
			index:    index,
			typ:      parameterType,
			optional: false,
		}

		definition.constructorArgs = append(definition.constructorArgs, arg)
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

func (d *Definition) Type() reflector.Type {
	return d.typ
}

func (d *Definition) Constructor() reflector.Function {
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

func (d *Definition) ConstructorArguments() []*ConstructorArgument {
	copyOfArgs := make([]*ConstructorArgument, 0)

	for _, arg := range d.constructorArgs {
		copyOfArgs = append(copyOfArgs, arg)
	}

	return copyOfArgs
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
		return fmt.Errorf("definition should not be nil")
	}

	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[definition.Name()]; exists {
		return fmt.Errorf("definition with name %s already exists", definition.Name())
	}

	r.definitionMap[definition.Name()] = definition

	return nil
}

func (r *ObjectDefinitionRegistry) Remove(name string) error {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[name]; !exists {
		return fmt.Errorf("no found definition with name %s", name)
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
	definitionList := r.List(filters...)

	if len(definitionList) > 1 {
		return nil, errors.New("definitions cannot be distinguished because too many matching found")
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

		if definition.Type().CanConvert(filterOpts.Type) {
			if reflector.IsStruct(definition.Type()) && reflector.IsStruct(filterOpts.Type) {
				if matchTypeName(definition.Type(), filterOpts.Type) {
					definitionList = append(definitionList, definition)
				}
			} else {
				definitionList = append(definitionList, definition)
			}
		} else if reflector.IsPointer(definition.Type()) && !reflector.IsPointer(filterOpts.Type) && !reflector.IsInterface(filterOpts.Type) {
			pointerType := reflector.ToPointer(definition.Type())

			if pointerType.Elem().CanConvert(filterOpts.Type) {
				if reflector.IsStruct(pointerType) && reflector.IsStruct(filterOpts.Type) {
					if matchTypeName(pointerType, filterOpts.Type) {
						definitionList = append(definitionList, definition)
					}
				} else {
					definitionList = append(definitionList, definition)
				}
			}
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
