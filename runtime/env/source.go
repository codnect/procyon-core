package env

import (
	"os"
	"strings"
)

type EnvironmentPropertySource struct {
	variables map[string]string
}

func NewPropertySource() *EnvironmentPropertySource {
	source := &EnvironmentPropertySource{
		variables: make(map[string]string, 0),
	}

	variables := os.Environ()

	for _, variable := range variables {
		index := strings.Index(variable, "=")

		if index != -1 {
			source.variables[variable[:index]] = variable[index+1:]
		}
	}

	return source
}

func (s *EnvironmentPropertySource) Name() string {
	return "systemEnvironment"
}

func (s *EnvironmentPropertySource) Source() any {
	copyOfVariables := make(map[string]string)
	for key, value := range s.variables {
		copyOfVariables[key] = value
	}

	return copyOfVariables
}

func (s *EnvironmentPropertySource) Property(name string) (any, bool) {
	propertyName, exists := s.checkPropertyName(strings.ToLower(name))

	if exists {
		if value, ok := s.variables[propertyName]; ok {
			return value, true
		}
	}

	propertyName, exists = s.checkPropertyName(strings.ToUpper(name))

	if exists {
		if value, ok := s.variables[propertyName]; ok {
			return value, true
		}
	}

	return nil, false
}

func (s *EnvironmentPropertySource) PropertyOrDefault(name string, defaultValue any) any {
	value, ok := s.Property(name)

	if !ok {
		return defaultValue
	}

	return value
}

func (s *EnvironmentPropertySource) ContainsProperty(name string) bool {
	_, exists := s.checkPropertyName(strings.ToUpper(name))
	if exists {
		return true
	}

	_, exists = s.checkPropertyName(strings.ToLower(name))
	if exists {
		return true
	}

	return false
}

func (s *EnvironmentPropertySource) PropertyNames() []string {
	keys := make([]string, 0, len(s.variables))

	for key, _ := range s.variables {
		keys = append(keys, key)
	}

	return keys
}

func (s *EnvironmentPropertySource) checkPropertyName(name string) (string, bool) {
	if s.contains(name) {
		return name, true
	}

	noHyphenPropertyName := strings.ReplaceAll(name, "-", "_")

	if name != noHyphenPropertyName && s.contains(noHyphenPropertyName) {
		return noHyphenPropertyName, true
	}

	noDotPropertyName := strings.ReplaceAll(name, ".", "_")

	if name != noDotPropertyName && s.contains(noDotPropertyName) {
		return noDotPropertyName, true
	}

	noHyphenAndNoDotName := strings.ReplaceAll(noDotPropertyName, "-", "_")

	if noDotPropertyName != noHyphenAndNoDotName && s.contains(noHyphenAndNoDotName) {
		return noHyphenAndNoDotName, true
	}

	return "", false
}

func (s *EnvironmentPropertySource) contains(name string) bool {
	if _, ok := s.variables[name]; ok {
		return true
	}

	return false
}
