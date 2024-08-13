package property

import "codnect.io/procyon-core/component"

func init() {
	// loaders
	component.Register(newYamlSourceLoader, component.Named("procyonYamlPropertySourceLoader"))
}
