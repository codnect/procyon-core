package runtime

import (
	"github.com/codnect/procyoncore/runtime/config"
	"github.com/codnect/procyoncore/runtime/env"
	"github.com/codnect/procyoncore/runtime/env/property"
	"strings"
)

type customizer struct {
	sourceLoaders []property.SourceLoader
}

func newCustomizer() *customizer {
	return &customizer{
		sourceLoaders: []property.SourceLoader{
			property.NewYamlPropertySourceLoader(),
		},
	}
}

func (c *customizer) CustomizeEnvironment(environment env.Environment) error {
	return c.importConfig(environment)
}

func (c *customizer) importConfig(environment env.Environment) error {
	importer := config.NewFileImporter(environment)

	defaultConfigs, err := importer.Load(environment.DefaultProfiles(), "resources")
	if err != nil {
		return err
	}

	sources := property.NewPropertySources()

	for _, defaultConfig := range defaultConfigs {
		sources.AddLast(defaultConfig.PropertySource())
	}

	activeProfiles := environment.ActiveProfiles()

	if len(activeProfiles) == 0 {
		resolver := property.NewSourcesResolver(sources)
		value, ok := resolver.Property("procyon.profiles.active")

		if ok {
			activeProfiles = strings.Split(strings.TrimSpace(value), ",")
		}
	}

	if len(activeProfiles) != 0 {
		err = environment.SetActiveProfiles(activeProfiles...)
		if err != nil {
			return err
		}

		err = c.loadActiveProfiles(importer, environment, sources, activeProfiles)
		if err != nil {
			return err
		}
	}

	c.mergeSources(environment, sources)
	return nil
}

func (c *customizer) loadActiveProfiles(importer config.Importer, environment env.Environment, propertySources *property.Sources, activeProfiles []string) error {
	configs, err := importer.Load(activeProfiles, "config")
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		propertySource := cfg.PropertySource()
		propertySources.AddFirst(propertySource)

		err = c.activateIncludeProfiles(importer, environment, propertySources, propertySource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *customizer) activateIncludeProfiles(importer config.Importer, environment env.Environment, propertySources *property.Sources, source property.Source) error {
	value, ok := source.Property("procyon.profiles.include")

	if ok {
		profiles := strings.Split(strings.TrimSpace(value.(string)), ",")

		for _, profile := range profiles {
			err := environment.AddActiveProfile(strings.TrimSpace(profile))
			if err != nil {
				return err
			}
		}

		err := c.loadActiveProfiles(importer, environment, propertySources, profiles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *customizer) mergeSources(environment env.Environment, sources *property.Sources) {
	for _, source := range sources.ToSlice() {
		environment.PropertySources().AddLast(source)
	}
}
