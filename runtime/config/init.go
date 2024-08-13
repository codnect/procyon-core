package config

import "codnect.io/procyon-core/component"

func init() {
	// resolvers
	component.Register(newDefaultDefaultResolver, component.Named("procyonDefaultConfigResourceResolver"))
	// loaders
	component.Register(newFileLoader, component.Named("procyonConfigFileLoader"))
	// other
	component.Register(newImporter, component.Named("procyonConfigImporter"))
}
