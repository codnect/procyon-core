package config

import (
	"context"
	"fmt"
	"github.com/codnect/procyoncore/runtime/property"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	FileName = "procyon"
)

type Resolver interface {
	Resolve(ctx context.Context, location string, profiles []string) ([]Resource, error)
}

type DefaultResolver struct {
	loaders    []property.SourceLoader
	configName string
}

func newDefaultResolver(loaders []property.SourceLoader) *DefaultResolver {
	return &DefaultResolver{
		loaders:    loaders,
		configName: FileName,
	}
}

func (r *DefaultResolver) Resolve(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	resources := make([]Resource, 0)
	if profiles == nil {
		resources = append(resources, r.getResources("", location)...)
		return resources, nil
	}

	for _, profile := range profiles {
		if profile == "default" {
			resources = append(resources, r.getResources("", location)...)
		} else {
			resources = append(resources, r.getResources(profile, location)...)
		}
	}

	return resources, nil
}

func (r *DefaultResolver) getResources(profile string, location string) []Resource {
	var (
		configFile fs.File
		resources  = make([]Resource, 0)
	)

	for _, loader := range r.loaders {
		extensions := loader.FileExtensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = filepath.Join(location, fmt.Sprintf("%s.%s", r.configName, extension))
			} else {
				filePath = filepath.Join(location, fmt.Sprintf("%s-%s.%s", r.configName, profile, extension))
			}

			if _, err := os.Stat(filePath); err == nil {
				configFile, err = os.Open(filePath)

				if err != nil {
					continue
				}

				resources = append(resources, newFileResource(filePath, configFile, loader))
			}
		}
	}

	return resources
}
