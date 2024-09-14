package property

import (
	"strconv"
	"strings"
)

type MapSource struct {
	name   string
	source map[string]interface{}
}

func NewMapSource(name string, source map[string]interface{}) *MapSource {
	if strings.TrimSpace(name) == "" {
		panic("cannot create map source with empty or blank name")
	}

	if source == nil {
		panic("nil source")
	}

	return &MapSource{
		name:   name,
		source: flatMap(source),
	}
}

func (m *MapSource) Name() string {
	return m.name
}

func (m *MapSource) Source() any {
	return m.source
}

func (m *MapSource) ContainsProperty(name string) bool {
	if _, exists := m.source[name]; exists {
		return true
	}

	return false
}

func (m *MapSource) Property(name string) (any, bool) {
	if value, exists := m.source[name]; exists {
		return value, true
	}

	return nil, false
}

func (m *MapSource) PropertyOrDefault(name string, defaultValue any) any {
	value, exists := m.Property(name)
	if !exists {
		return defaultValue
	}

	return value
}

func (m *MapSource) PropertyNames() []string {
	names := make([]string, 0)

	for name, _ := range m.source {
		names = append(names, name)
	}

	return names
}

func flatMap(m map[string]interface{}) map[string]interface{} {
	flattenMap := map[string]interface{}{}

	for key, value := range m {
		switch child := value.(type) {
		case map[string]interface{}:
			nm := flatMap(child)

			for nk, nv := range nm {
				flattenMap[key+"."+nk] = nv
			}
		case []interface{}:
			for i := 0; i < len(child); i++ {
				flattenMap[key+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			flattenMap[key] = value
		}
	}

	return flattenMap
}
