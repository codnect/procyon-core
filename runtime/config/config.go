package config

import (
	"github.com/codnect/procyoncore/runtime/property"
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
