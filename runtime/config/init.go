package config

import "github.com/codnect/procyoncore/component"

func init() {
	// resolvers
	component.Register(newDefaultResolver, component.Named("procyonDefaultConfigResolver"))
	// loaders
	component.Register(newFileLoader, component.Named("procyonConfigFileLoader"))
	// other
	component.Register(newImporter, component.Named("procyonConfigImporter"))
}
