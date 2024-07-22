package component

import "sync"

var (
	components   = make([]*Component, 0)
	muComponents = sync.RWMutex{}
)

type Component struct {
	definition *Definition
	conditions []Condition
}

func newComponent(constructor Constructor) *Component {
	definition, err := MakeDefinition(constructor)

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

	for _, condition := range copyOfConditions {
		copyOfConditions = append(copyOfConditions, condition)
	}

	return copyOfConditions
}

type Registration struct {
	component *Component
}

func newRegistration(constructor Constructor) Registration {
	return Registration{
		component: newComponent(constructor),
	}
}

func (r Registration) Named(name string) Registration {
	WithName(name)(r.component.definition)
	return r
}

func (r Registration) Scoped(scope string) Registration {
	WithScope(scope)(r.component.definition)
	return r
}

func (r Registration) Primary() Registration {
	WithPrimary()(r.component.definition)
	return r
}

func (r Registration) Prioritized(priority int) Registration {
	WithPriority(priority)(r.component.definition)
	return r
}

func (r Registration) Inject(argumentIndex int) Registration {

	return r
}

func (r Registration) ConditionalOn(condition Condition) Registration {
	if condition != nil {
		r.component.conditions = append(r.component.conditions, condition)
	}

	return r
}

func Register(constructor Constructor) Registration {
	defer muComponents.Unlock()
	muComponents.Lock()
	registration := newRegistration(constructor)

	components = append(components, registration.component)

	return registration
}

func List() []*Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	copyOfComponents := make([]*Component, 0)
	for _, component := range components {
		copyOfComponents = append(copyOfComponents, component)
	}

	return copyOfComponents
}
