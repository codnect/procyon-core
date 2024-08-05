package runtime

import (
	"context"
	"github.com/codnect/procyoncore/runtime/config"
	"github.com/codnect/procyoncore/runtime/property"
	"strings"
)

type customizer struct {
	loaders  []property.SourceLoader
	importer *config.Importer
}

func newCustomizer(loaders []property.SourceLoader, importer *config.Importer) *customizer {
	return &customizer{
		loaders:  loaders,
		importer: importer,
	}
}

func (c *customizer) CustomizeEnvironment(environment Environment) error {
	return c.importConfig(environment)
}

func (c *customizer) importConfig(environment Environment) error {
	defaultConfigs, err := c.importer.LoadConfigs(context.Background(), "resources", environment.DefaultProfiles())

	if err != nil {
		return err
	}

	sourceList := property.NewSourceList()

	for _, defaultConfig := range defaultConfigs {
		sourceList.AddLast(defaultConfig.PropertySource())
	}

	activeProfiles := environment.ActiveProfiles()

	if len(activeProfiles) == 0 {
		resolver := property.NewMultiSourceResolver(sourceList)
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

		err = c.loadActiveProfiles(environment, sourceList, activeProfiles)
		if err != nil {
			return err
		}
	}

	c.mergeSources(environment, sourceList)
	return nil
}

func (c *customizer) loadActiveProfiles(environment Environment, sourceList *property.SourceList, activeProfiles []string) error {
	configs, err := c.importer.LoadConfigs(context.Background(), "config", activeProfiles)
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

func (c *customizer) activateIncludeProfiles(environment Environment, sourceList *property.SourceList, source property.Source) error {
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

func (c *customizer) mergeSources(environment Environment, sourceList *property.SourceList) {
	for _, source := range sourceList.ToSlice() {
		environment.PropertySources().AddLast(source)
	}
}
