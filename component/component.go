package component

import (
	"codnect.io/procyon-core/component/filter"
	"context"
	"fmt"
	"sync"
)

var (
	components   = make(map[string]*Component)
	muComponents = sync.RWMutex{}
)

type Component struct {
	definition *Definition
	conditions []Condition
}

type Initialization interface {
	DoInit(ctx context.Context) error
}

func createComponent(constructorFunc ConstructorFunc, options ...Option) *Component {
	definition, err := MakeDefinition(constructorFunc, options...)

	if err != nil {
		panic(err)
	}

	return &Component{
		definition: definition,
		conditions: make([]Condition, 0),
	}
}

func (c *Component) Definition() *Definition {
	return c.definition
}

func (c *Component) Conditions() []Condition {
	copyOfConditions := make([]Condition, 0)

	for _, condition := range c.conditions {
		copyOfConditions = append(copyOfConditions, condition)
	}

	return copyOfConditions
}

type Registration struct {
	component *Component
}

func (r Registration) ConditionalOn(condition Condition) Registration {
	if condition != nil {
		r.component.conditions = append(r.component.conditions, condition)
	}

	return r
}

func Register(constructorFunc ConstructorFunc, options ...Option) Registration {
	defer muComponents.Unlock()
	muComponents.Lock()

	component := createComponent(constructorFunc, options...)
	componentName := component.Definition().Name()

	if _, exists := components[componentName]; exists {
		panic(fmt.Sprintf("component with name '%s' already exists", componentName))
	}

	components[componentName] = component

	return Registration{
		component: component,
	}
}

func List(filters ...filter.Filter) []*Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	filterOpts := filter.Of(filters...)
	componentList := make([]*Component, 0)

	for _, component := range components {
		definition := component.Definition()

		if filterOpts.Name != "" && filterOpts.Name != component.Definition().Name() {
			continue
		}

		if filterOpts.Type == nil {
			componentList = append(componentList, component)
			continue
		}

		if matchTypes(definition.Type(), filterOpts.Type) {
			componentList = append(componentList, component)
		}
	}

	return componentList
}
