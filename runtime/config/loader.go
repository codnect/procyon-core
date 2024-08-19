package config

import (
	"codnect.io/reflector"
	"context"
	"errors"
	"fmt"
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
		return nil, errors.New("resource cannot be nil")
	}

	if fileResource, ok := resource.(*FileResource); ok {
		loader := fileResource.Loader()
		source, err := loader.Load(fileResource.Name(), fileResource.File())

		if err != nil {
			return nil, err
		}

		return New(source), err
	}

	return nil, fmt.Errorf("resource %s is not supported", reflector.TypeOfAny(resource).Name())
}
