package config

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Loader interface {
	IsLoadable(resource Resource) bool
	LoadConfig(ctx context.Context, resource Resource) (*Config, error)
}

type FileLoader struct {
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (l *FileLoader) IsLoadable(resource Resource) bool {
	_, canConvert := resource.(*FileResource)
	return canConvert
}

func (l *FileLoader) LoadConfig(ctx context.Context, resource Resource) (*Config, error) {
	if resource == nil {
		return nil, errors.New("nil context")
	}

	if resource == nil {
		return nil, errors.New("nil resource")
	}

	if fileResource, ok := resource.(*FileResource); ok {
		loader := fileResource.Loader()
		source, err := loader.Load(fileResource.Name(), fileResource.File())

		if err != nil {
			return nil, err
		}

		return New(source), err
	}

	return nil, fmt.Errorf("resource '%s' is not supported", reflect.TypeOf(resource).Name())
}
