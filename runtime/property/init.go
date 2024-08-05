package property

import "github.com/codnect/procyoncore/component"

func init() {
	// loaders
	component.Register(newYamlSourceLoader, component.Named("procyonYamlPropertySourceLoader"))
}
