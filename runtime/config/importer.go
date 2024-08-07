package config

import (
	"codnect.io/reflector"
	"context"
	"fmt"
)

type Importer struct {
	resolvers []ResourceResolver
	loaders   []Loader
}

func newImporter(resolvers []ResourceResolver, loaders []Loader) *Importer {
	return &Importer{
		resolvers: resolvers,
		loaders:   loaders,
	}
}

func (i *Importer) Import(ctx context.Context, location string, profiles []string) ([]*Config, error) {
	resources, err := i.resolve(ctx, location, profiles)
	if err != nil {
		return nil, err
	}

	return i.load(ctx, resources)
}

func (i *Importer) resolve(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, resolver := range i.resolvers {
		resolved, err := resolver.ResolveResources(ctx, location, profiles)

		if err != nil {
			return nil, err
		}

		resources = append(resources, resolved...)
	}

	return resources, nil
}

func (i *Importer) load(ctx context.Context, resources []Resource) ([]*Config, error) {
	loaded := make([]*Config, 0)

	for _, resource := range resources {
		loader, err := i.findLoader(resource)

		if err != nil {
			return nil, err
		}

		var data *Config
		data, err = loader.LoadConfig(ctx, resource)

		loaded = append(loaded, data)
	}

	return loaded, nil
}

func (i *Importer) findLoader(resource Resource) (Loader, error) {
	var result Loader
	for _, loader := range i.loaders {
		if loader.IsLoadable(resource) {

			if result != nil {
				return nil, fmt.Errorf("multiple loaders found for resource '%s'", reflector.TypeOfAny(resource).Name())
			}

			result = loader
		}
	}

	if result == nil {
		return nil, fmt.Errorf("no loader found for resource '%s'", reflector.TypeOfAny(resource).Name())
	}

	return result, nil
}
