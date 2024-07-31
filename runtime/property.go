package runtime

import (
	"strings"
)

const (
	NonOptionArgs = "nonOptionArgs"
)

type PropertySource struct {
	args *Arguments
}

func NewPropertySource(args *Arguments) *PropertySource {
	return &PropertySource{
		args: args,
	}
}

func (s *PropertySource) Name() string {
	return "commandLineArgs"
}

func (s *PropertySource) Source() any {
	return s.args
}

func (s *PropertySource) ContainsProperty(name string) bool {
	if NonOptionArgs == name {
		return len(s.args.NonOptionArgs()) != 0
	}

	return s.args.ContainsOption(name)
}

func (s *PropertySource) Property(name string) (any, bool) {
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

func (s *PropertySource) PropertyOrDefault(name string, defaultValue any) any {
	val, ok := s.Property(name)
	if !ok {
		return defaultValue
	}

	return val
}

func (s *PropertySource) PropertyNames() []string {
	return s.args.OptionNames()
}

func (s *PropertySource) OptionValues(name string) []string {
	return s.args.OptionValues(name)
}

func (s *PropertySource) NonOptionArgs() []string {
	return s.args.NonOptionArgs()
}
