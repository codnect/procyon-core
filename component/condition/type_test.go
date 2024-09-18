package condition

/*
type AnyType struct {
}

func anyConstructorFunction() AnyType {
	return AnyType{}
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewContainer()
	err := container.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyDefinitionWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewContainer()

	definition, err := container.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	container := component.NewContainer()
	err := container.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyDefinitionWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	container := component.NewContainer()

	definition, err := container.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onMissingTypeCondition := OnMissingObject("anyObjectName")
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}
*/
