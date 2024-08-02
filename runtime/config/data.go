package config

import (
	"github.com/codnect/procyoncore/runtime/property"
)

type Data struct {
	source property.Source
}

func NewData(source property.Source) *Data {
	return &Data{
		source,
	}
}

func (d *Data) PropertySource() property.Source {
	return d.source
}
