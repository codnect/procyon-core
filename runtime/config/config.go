package config

import (
	"codnect.io/procyon-core/runtime/property"
)

type Config struct {
	source property.Source
}

func New(source property.Source) *Config {
	return &Config{
		source,
	}
}

func (d *Config) PropertySource() property.Source {
	return d.source
}
