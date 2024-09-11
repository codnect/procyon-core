package component

import (
	"codnect.io/procyon-core/component/filter"
	"codnect.io/reflector"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingletonObjectRegistry_RegisterShouldRegisterObjectSuccessfully(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)
	assert.Contains(t, registry.singletonObjects, "anyObjectName")
	assert.Contains(t, registry.typesOfSingletonObjects, "anyObjectName")
	assert.Equal(t, reflector.TypeOf[*AnyType]().ReflectType(),
		registry.typesOfSingletonObjects["anyObjectName"].ReflectType())
}

func TestSingletonObjectRegistry_RegisterShouldReturnErrorIfComponentWithSameNameIsAlreadyRegistered(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)

	err = registry.Register("anyObjectName", &AnyType{})
	assert.Equal(t, "object with name 'anyObjectName' already exists", err.Error())
}

func TestSingletonObjectRegistry_ContainsShouldReturnTrueIfComponentExists(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)
	assert.True(t, registry.Contains("anyObjectName"))
}

func TestSingletonObjectRegistry_ContainsShouldReturnFalseIfComponentDoesNotExist(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	assert.False(t, registry.Contains("anyObjectName"))
}

func TestSingletonObjectRegistry_Find(t *testing.T) {

	type args struct {
		filter []filter.Filter
	}

	type fields struct {
		singletonObjects        map[string]any
		typesOfSingletonObjects map[string]reflector.Type
	}

	anyObject := &AnyType{}
	anotherObject := &AnotherType{}
	anyObjectType := reflector.TypeOfAny(anyObject)
	anotherObjectType := reflector.TypeOfAny(anotherObject)

	testCases := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr string
	}{
		{
			name: "ShouldReturnObjectWithoutFiltersIfThereIsOnlyOneObject",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName": anyObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName": anyObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want: anyObject,
		},
		{
			name: "ShouldReturnErrorWithoutFiltersIfThereAreManyObjects",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:    nil,
			wantErr: "cannot distinguish objects because too many matching found",
		},
		{
			name: "ShouldReturnObjectWithByNameFilterIfObjectWithNameExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want: anyObject,
		},
		{
			name: "ShouldReturnErrorWithByNameFilterIfObjectWithNameDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:    nil,
			wantErr: "not found object with name 'anyObjectName'",
		},
		{
			name: "ShouldReturnObjectWithByTypeOfFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: anyObject,
		},
		{
			name: "ShouldReturnErrorWithByTypeOfFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:    nil,
			wantErr: "not found object with type '*AnyType'",
		},
		{
			name: "ShouldReturnObjectWithByTypeOfFilterIfThereIsOnlyOneObjectImplementingInterface",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnotherInterface](),
				},
			},
			want: anotherObject,
		},
		{
			name: "ShouldReturnErrorWithByTypeOfFilterIfThereIsMoreThanOneObjectImplementingInterface",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnyInterface](),
				},
			},
			want:    nil,
			wantErr: "cannot distinguish objects because too many matching found",
		},
		{
			name: "ShouldReturnObjectWithByTypeFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want: anyObject,
		},
		{
			name: "ShouldReturnErrorWithByTypeFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:    nil,
			wantErr: "not found object with type '*AnyType'",
		},
		{
			name: "ShouldReturnObjectWithAllFiltersIfObjectExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: anyObject,
		},
		{
			name: "ShouldReturnErrorWithAllFiltersIfObjectDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:    nil,
			wantErr: "not found object with name 'anyObjectName' and type '*AnyType'",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry := NewSingletonObjectRegistry()
			registry.singletonObjects = testCase.fields.singletonObjects
			registry.typesOfSingletonObjects = testCase.fields.typesOfSingletonObjects

			got, err := registry.Find(testCase.args.filter...)
			if testCase.wantErr != "" {
				if err != nil {
					assert.Equal(t, testCase.wantErr, err.Error(), "failed for test case '%s'", testCase.name)
				} else {
					assert.Nil(t, err, "want err '%s' but got nil", testCase.wantErr)
				}
			}

			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestSingletonObjectRegistry_FindFirst(t *testing.T) {
	type args struct {
		filter []filter.Filter
	}

	type fields struct {
		singletonObjects        map[string]any
		typesOfSingletonObjects map[string]reflector.Type
	}

	anyObject := &AnyType{}
	anotherObject := &AnotherType{}
	anyObjectType := reflector.TypeOfAny(anyObject)
	anotherObjectType := reflector.TypeOfAny(anotherObject)

	testCases := []struct {
		name     string
		fields   fields
		args     args
		want     any
		wantIn   []any
		wantBool bool
	}{
		{
			name: "ShouldReturnObjectAndTrueWithoutFiltersIfThereIsOnlyOneObject",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName": anyObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName": anyObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:     anyObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnObjectAndTrueWithoutFiltersIfThereAreManyObjects",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			wantIn:   []any{anyObject, anotherObject},
			wantBool: true,
		},
		{
			name: "ShouldReturnObjectAndTrueWithByNameFilterIfObjectWithNameExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:     anyObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnNilAndFalseWithByNameFilterIfObjectWithNameDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnObjectAndTrueWithByTypeOfFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     anyObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnNilAndFalseWithByTypeOfFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnObjectAndTrueWithByTypeOfFilterIfThereIsOnlyOneObjectImplementingInterface",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnotherInterface](),
				},
			},
			want:     anotherObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnObjectAndTrueWithByTypeFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:     anyObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnNilAndFalseWithByTypeFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnObjectAndTrueWithAllFiltersIfObjectExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     anyObject,
			wantBool: true,
		},
		{
			name: "ShouldReturnNilAndFalseWithAllFiltersIfObjectDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     nil,
			wantBool: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry := NewSingletonObjectRegistry()
			registry.singletonObjects = testCase.fields.singletonObjects
			registry.typesOfSingletonObjects = testCase.fields.typesOfSingletonObjects

			got, ok := registry.FindFirst(testCase.args.filter...)
			assert.Equal(t, testCase.wantBool, ok)
			if len(testCase.wantIn) != 0 {
				found := false
				for _, want := range testCase.wantIn {
					if want == got {
						found = true
						break
					}
				}

				assert.True(t, found, "not found any matching object in wantIn %v", testCase.want)
			} else {
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestSingletonObjectRegistry_OrElseCreateShouldCreateAndReturnObjectIfObjectDoesNotExist(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	anyObject := &AnyType{}
	methodCalled := false

	got, err := registry.OrElseCreate("anyObjectName", func(ctx context.Context) (any, error) {
		methodCalled = true
		return anyObject, nil
	})

	assert.Nil(t, err)
	assert.True(t, methodCalled)
	assert.Equal(t, anyObject, got)
}

func TestSingletonObjectRegistry_OrElseCreateShouldReturnErrorIfObjectDoesNotExistAndProviderReturnError(t *testing.T) {
	registry := NewSingletonObjectRegistry()

	anyError := errors.New("anyError")
	got, err := registry.OrElseCreate("anyObjectName", func(ctx context.Context) (any, error) {
		return nil, anyError
	})

	assert.Equal(t, anyError.Error(), err.Error())
	assert.Nil(t, got)
}

func TestSingletonObjectRegistry_OrElseCreateShouldReturnObjectIfObjectAlreadyExists(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	anyObject := &AnyType{}
	err := registry.Register("anyObjectName", anyObject)
	assert.Nil(t, err)

	methodCalled := false
	var got any
	got, err = registry.OrElseCreate("anyObjectName", func(ctx context.Context) (any, error) {
		methodCalled = true
		return &AnyType{}, nil
	})

	assert.Nil(t, err)
	assert.False(t, methodCalled)
	assert.Equal(t, anyObject, got)
}

func TestSingletonObjectRegistry_OrElseCreateShouldReturnErrorIfObjectWithSameNameIsAlreadyInPreparation(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	registry.singletonObjectsInPreparation["anyObjectName"] = struct{}{}

	got, err := registry.OrElseCreate("anyObjectName", func(ctx context.Context) (any, error) {
		return &AnyType{}, nil
	})

	assert.Equal(t, "object with name 'anyObjectName' is currently in preparation, maybe it has got circular dependency cycle", err.Error())
	assert.Nil(t, got)
}

func TestSingletonObjectRegistry_CountShouldReturnCountOfObjects(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)

	assert.Equal(t, 1, len(registry.singletonObjects))
	assert.Equal(t, 1, len(registry.typesOfSingletonObjects))
	assert.Equal(t, 1, registry.Count())
}

func TestSingletonObjectRegistry_NamesShouldReturnListOfObjectNames(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)
	err = registry.Register("anotherObjectName", &AnotherType{})
	assert.Nil(t, err)

	names := registry.Names()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "anyObjectName")
	assert.Contains(t, names, "anotherObjectName")
}

func TestSingletonObjectRegistry_RemoveShouldDeleteObjectFromRegistryIfObjectWithNameExists(t *testing.T) {
	registry := NewSingletonObjectRegistry()
	err := registry.Register("anyObjectName", &AnyType{})
	assert.Nil(t, err)

	err = registry.Remove("anyObjectName")
	assert.Nil(t, err)

	assert.Zero(t, len(registry.singletonObjects))
	assert.Zero(t, len(registry.typesOfSingletonObjects))
	assert.Zero(t, registry.Count())
}

func TestSingletonObjectRegistry_RemoveShouldReturnErrorIfObjectWithNameDoesNotExist(t *testing.T) {
	registry := NewSingletonObjectRegistry()

	err := registry.Remove("anyObjectName")
	assert.Equal(t, "no found object with name 'anyObjectName'", err.Error())
}

func TestSingletonObjectRegistry_List(t *testing.T) {
	type args struct {
		filter []filter.Filter
	}

	type fields struct {
		singletonObjects        map[string]any
		typesOfSingletonObjects map[string]reflector.Type
	}

	anyObject := &AnyType{}
	anotherObject := &AnotherType{}
	anyObjectType := reflector.TypeOfAny(anyObject)
	anotherObjectType := reflector.TypeOfAny(anotherObject)

	testCases := []struct {
		name   string
		fields fields
		args   args
		want   []any
	}{
		{
			name: "ShouldReturnObjectsWithoutFiltersIfThereIsOnlyOneObject",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName": anyObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName": anyObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want: []any{anyObject},
		},
		{
			name: "ShouldReturnObjectsWithoutFiltersIfThereAreManyObjects",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want: []any{anyObject, anotherObject},
		},
		{
			name: "ShouldReturnObjectWithByNameFilterIfObjectWithNameExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want: []any{anyObject},
		},
		{
			name: "ShouldReturnEmptySliceWithByNameFilterIfObjectWithNameDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want: []any{},
		},
		{
			name: "ShouldReturnObjectWithByTypeOfFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: []any{anyObject},
		},
		{
			name: "ShouldReturnEmptySliceWithByTypeOfFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: []any{},
		},
		{
			name: "ShouldReturnObjectWithByTypeOfFilterIfThereIsOnlyOneObjectImplementingInterface",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnotherInterface](),
				},
			},
			want: []any{anotherObject},
		},
		{
			name: "ShouldReturnObjectsWithByTypeOfFilterIfThereIsMoreThanOneObjectImplementingInterface",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnyInterface](),
				},
			},
			want: []any{anyObject, anotherObject},
		},
		{
			name: "ShouldReturnObjectWithByTypeFilterIfObjectWithTypeExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want: []any{anyObject},
		},
		{
			name: "ShouldReturnEmptySliceWithByTypeFilterIfObjectWithTypeDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want: []any{},
		},
		{
			name: "ShouldReturnObjectWithAllFiltersIfObjectExists",
			fields: fields{
				singletonObjects: map[string]any{
					"anyObjectName":     anyObject,
					"anotherObjectName": anotherObject,
				},
				typesOfSingletonObjects: map[string]reflector.Type{
					"anyObjectName":     anyObjectType,
					"anotherObjectName": anotherObjectType,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: []any{anyObject},
		},
		{
			name: "ShouldReturnEmptySliceWithAllFiltersIfObjectDoesNotExist",
			fields: fields{
				singletonObjects:        map[string]any{},
				typesOfSingletonObjects: map[string]reflector.Type{},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: []any{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry := NewSingletonObjectRegistry()
			registry.singletonObjects = testCase.fields.singletonObjects
			registry.typesOfSingletonObjects = testCase.fields.typesOfSingletonObjects

			got := registry.List(testCase.args.filter...)
			assert.Equal(t, len(testCase.want), len(got))
			for _, w := range testCase.want {
				assert.Contains(t, got, w)
			}
		})
	}
}
