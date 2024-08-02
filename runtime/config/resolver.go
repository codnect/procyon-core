package config

import (
	"fmt"
	"github.com/codnect/procyoncore/runtime"
	"github.com/codnect/procyoncore/runtime/property"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type LocationResolver interface {
	Resolve(location string, profiles []string) ([]Resource, error)
}

type FileLocationResolver struct {
	environment   runtime.Environment
	sourceLoaders []property.SourceLoader
	configName    string
}

func NewFileLocationResolver(environment runtime.Environment, sourceLoaders []property.SourceLoader) *FileLocationResolver {
	if environment == nil {
		panic("environment cannot be nil")
	}

	if len(sourceLoaders) == 0 {
		panic("sourceLoaders cannot be empty")
	}

	resolver := &FileLocationResolver{
		environment:   environment,
		sourceLoaders: sourceLoaders,
	}

	configNameProperty := environment.PropertyResolver().PropertyOrDefault("procyon.config.name", "procyon")
	resolver.configName = strings.TrimSpace(configNameProperty.(string))

	if resolver.configName == "" {
		panic("configName cannot be empty or blank")
	}

	return resolver
}

func (r *FileLocationResolver) Resolve(location string, profiles []string) ([]Resource, error) {
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

func (r *FileLocationResolver) getResources(profile string, location string) []Resource {
	resources := make([]Resource, 0)
	var configFile fs.File

	for _, loader := range r.sourceLoaders {
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

				resources = append(resources, NewFileResource(filePath, configFile, loader))
			}
		}
	}

	return resources
}
