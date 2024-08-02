package config

import (
	"github.com/codnect/procyoncore/runtime"
	"github.com/codnect/procyoncore/runtime/property"
)

type FileImporter struct {
	resolver *FileLocationResolver
}

func NewFileImporter(environment runtime.Environment) *FileImporter {
	resolver := NewFileLocationResolver(environment, []property.SourceLoader{
		property.NewYamlSourceLoader(),
	})

	return &FileImporter{
		resolver,
	}
}

func (i *FileImporter) Load(location string, profiles []string) ([]*Data, error) {
	resources, err := i.resolver.Resolve(location, profiles)
	if err != nil {
		return nil, err
	}

	return i.loadResources(resources)
}

func (i *FileImporter) loadResources(resources []Resource) ([]*Data, error) {
	configs := make([]*Data, 0)

	for _, resource := range resources {
		fileResource, canConvert := resource.(*FileResource)
		if !canConvert {
			continue
		}

		loader := resource.Loader()
		source, err := loader.Load(fileResource.Location(), fileResource.File())

		if err != nil {
			return nil, err
		}

		configs = append(configs, NewData(source))
	}

	return configs, nil
}
