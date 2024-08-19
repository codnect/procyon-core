package property

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type SourceLoader interface {
	FileExtensions() []string
	Load(name string, reader io.Reader) (Source, error)
}

type YamlSourceLoader struct {
}

func NewYamlSourceLoader() *YamlSourceLoader {
	return &YamlSourceLoader{}
}

func (l *YamlSourceLoader) FileExtensions() []string {
	return []string{"yml", "yaml"}
}

func (l *YamlSourceLoader) Load(name string, reader io.Reader) (Source, error) {
	loaded := make(map[string]interface{})

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}

	return NewMapSource(name, loaded), nil
}
