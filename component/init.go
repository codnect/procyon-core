package component

type InitializationContext struct {
	definition *Definition
	object     any
}

func newInitContext(definition *Definition, object any) InitializationContext {
	return InitializationContext{
		definition: definition,
		object:     object,
	}
}

func (c InitializationContext) Definition() *Definition {
	return c.definition
}

func (c InitializationContext) Object() any {
	return c.object
}

type BeforeInitialization interface {
	BeforeInit(ctx InitializationContext) (any, error)
}

type AfterInitialization interface {
	AfterInit(ctx InitializationContext) (any, error)
}
