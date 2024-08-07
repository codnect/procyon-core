package runtime

import (
	"context"
	"github.com/codnect/procyoncore/runtime/config"
	"github.com/codnect/procyoncore/runtime/property"
	"strings"
)

type defaultCustomizer struct {
	loaders  []property.SourceLoader
	importer *config.Importer
}

func newDefaultCustomizer(loaders []property.SourceLoader, importer *config.Importer) *defaultCustomizer {
	return &defaultCustomizer{
		loaders:  loaders,
		importer: importer,
	}
}

func (c *defaultCustomizer) CustomizeEnvironment(environment Environment) error {
	return c.importConfig(environment)
}

func (c *defaultCustomizer) importConfig(environment Environment) error {
	defaultConfigs, err := c.importer.Import(context.Background(), "resources", environment.DefaultProfiles())

	if err != nil {
		return err
	}

	sources := property.NewSources()

	for _, defaultConfig := range defaultConfigs {
		sources.AddLast(defaultConfig.PropertySource())
	}

	activeProfiles := environment.ActiveProfiles()

	if len(activeProfiles) == 0 {
		resolver := property.NewSourcesResolver(sources.ToSlice()...)
		value, ok := resolver.Property("procyon.profiles.active")

		if ok {
			activeProfiles = strings.Split(strings.TrimSpace(value.(string)), ",")
		}
	}

	if len(activeProfiles) != 0 {
		err = environment.SetActiveProfiles(activeProfiles...)
		if err != nil {
			return err
		}

		err = c.loadActiveProfiles(environment, sources, activeProfiles)
		if err != nil {
			return err
		}
	}

	c.mergeSources(environment, sources)
	return nil
}

func (c *defaultCustomizer) loadActiveProfiles(environment Environment, sourceList *property.Sources, activeProfiles []string) error {
	configs, err := c.importer.Import(context.Background(), "config", activeProfiles)
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		propertySource := cfg.PropertySource()
		sourceList.AddFirst(propertySource)

		err = c.activateIncludeProfiles(environment, sourceList, propertySource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCustomizer) activateIncludeProfiles(environment Environment, sourceList *property.Sources, source property.Source) error {
	value, ok := source.Property("procyon.profiles.include")

	if ok {
		profiles := strings.Split(strings.TrimSpace(value.(string)), ",")

		for _, profile := range profiles {
			err := environment.AddActiveProfile(strings.TrimSpace(profile))
			if err != nil {
				return err
			}
		}

		err := c.loadActiveProfiles(environment, sourceList, profiles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCustomizer) mergeSources(environment Environment, sourceList *property.Sources) {
	for _, source := range sourceList.ToSlice() {
		environment.PropertySources().AddLast(source)
	}
}
