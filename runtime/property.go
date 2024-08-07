package runtime

import (
	"os"
	"strings"
)

const (
	NonOptionArgs = "nonOptionArgs"
)

type ArgumentsSource struct {
	args *Arguments
}

func NewArgumentsSource(args *Arguments) *ArgumentsSource {
	return &ArgumentsSource{
		args: args,
	}
}

func (s *ArgumentsSource) Name() string {
	return "commandLineArgs"
}

func (s *ArgumentsSource) Source() any {
	return s.args
}

func (s *ArgumentsSource) ContainsProperty(name string) bool {
	if NonOptionArgs == name {
		return len(s.args.NonOptionArgs()) != 0
	}

	return s.args.ContainsOption(name)
}

func (s *ArgumentsSource) Property(name string) (any, bool) {
	if NonOptionArgs == name {
		nonOptValues := s.args.NonOptionArgs()

		if nonOptValues != nil {
			return strings.Join(nonOptValues, ","), true
		}

		return nil, false
	}

	optValues := s.args.OptionValues(name)

	if optValues != nil {
		return strings.Join(optValues, ","), true
	}

	return nil, false
}

func (s *ArgumentsSource) PropertyOrDefault(name string, defaultValue any) any {
	val, ok := s.Property(name)
	if !ok {
		return defaultValue
	}

	return val
}

func (s *ArgumentsSource) PropertyNames() []string {
	return s.args.OptionNames()
}

func (s *ArgumentsSource) OptionValues(name string) []string {
	return s.args.OptionValues(name)
}

func (s *ArgumentsSource) NonOptionArgs() []string {
	return s.args.NonOptionArgs()
}

type EnvironmentSource struct {
	variables map[string]string
}

func NewEnvironmentSource() *EnvironmentSource {
	source := &EnvironmentSource{
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

func (s *EnvironmentSource) Name() string {
	return "systemEnvironment"
}

func (s *EnvironmentSource) Source() any {
	copyOfVariables := make(map[string]string)
	for key, value := range s.variables {
		copyOfVariables[key] = value
	}

	return copyOfVariables
}

func (s *EnvironmentSource) Property(name string) (any, bool) {
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

func (s *EnvironmentSource) PropertyOrDefault(name string, defaultValue any) any {
	value, ok := s.Property(name)

	if !ok {
		return defaultValue
	}

	return value
}

func (s *EnvironmentSource) ContainsProperty(name string) bool {
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

func (s *EnvironmentSource) PropertyNames() []string {
	keys := make([]string, 0, len(s.variables))

	for key, _ := range s.variables {
		keys = append(keys, key)
	}

	return keys
}

func (s *EnvironmentSource) checkPropertyName(name string) (string, bool) {
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

func (s *EnvironmentSource) contains(name string) bool {
	if _, ok := s.variables[name]; ok {
		return true
	}

	return false
}
