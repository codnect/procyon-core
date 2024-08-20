package core

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/config"
	"codnect.io/procyon-core/runtime/property"
)

type Module struct {
}

func (m Module) InitModule() error {
	// core
	component.Register(newConfigContextConfigurer, component.Named("procyonConfigContextConfigurer"))
	// runtime/config
	component.Register(config.NewDefaultResourceResolver, component.Named("procyonDefaultConfigResourceResolver"))
	component.Register(config.NewFileLoader, component.Named("procyonConfigFileLoader"))
	component.Register(config.NewImporter, component.Named("procyonConfigImporter"))
	// runtime/property
	component.Register(property.NewYamlSourceLoader, component.Named("procyonYamlPropertySourceLoader"))
	// runtime
	component.Register(runtime.NewServerProperties, component.Scoped(component.PrototypeScope))
	component.Register(runtime.NewLifecycleProperties, component.Scoped(component.PrototypeScope))
	return nil
}
