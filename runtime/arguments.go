package runtime

type Arguments interface {
	OptionNames() []string
	ContainsOption(name string) bool
	NonOptionArgs() []string
	OptionValues(name string) []string
}
