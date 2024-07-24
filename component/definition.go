package component

import (
	"codnect.io/reflector"
	"fmt"
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
	Register(def *Definition) error
	Remove(name string) error
	Contains(name string) bool
	Find(name string) (*Definition, bool)
	List() []*Definition
	Names() []string
	NamesByType(requiredType reflector.Type) []string
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

func (r *ObjectDefinitionRegistry) Register(def *Definition) error {
	if def == nil {
		return fmt.Errorf("definition should not be nil")
	}

	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[def.Name()]; exists {
		return fmt.Errorf("definition with name %s already exists", def.Name())
	}

	r.definitionMap[def.Name()] = def

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

func (r *ObjectDefinitionRegistry) Find(name string) (*Definition, bool) {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if def, exists := r.definitionMap[name]; exists {
		return def, true
	}

	return nil, false
}

func (r *ObjectDefinitionRegistry) List() []*Definition {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	defs := make([]*Definition, 0)
	for _, def := range r.definitionMap {
		defs = append(defs, def)
	}

	return defs
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

func (r *ObjectDefinitionRegistry) NamesByType(requiredType reflector.Type) []string {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	names := make([]string, 0)

	if requiredType == nil {
		return names
	}

	for name, def := range r.definitionMap {

		instanceType := def.Type()

		if instanceType.CanConvert(requiredType) {
			names = append(names, name)
		} else if reflector.IsPointer(instanceType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
			ptrType := reflector.ToPointer(instanceType)

			if ptrType.Elem().CanConvert(requiredType) {
				names = append(names, name)
			}
		}

	}

	return names
}

func (r *ObjectDefinitionRegistry) Count() int {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	return len(r.definitionMap)
}
