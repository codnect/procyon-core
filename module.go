package core

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/condition"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/config"
	"codnect.io/procyon-core/runtime/event"
	"codnect.io/procyon-core/runtime/property"
)

type Module struct {
}

func (m Module) InitModule() error {
	// core
	component.Register(newConfigContextConfigurer, component.WithName("procyonConfigContextConfigurer"))
	// runtime/event
	component.Register(event.NewSimpleMulticaster, component.WithName("procyonEventMulticaster"),
		component.WithCondition(condition.OnMissingType[event.Multicaster]()),
	)
	// runtime/config
	component.Register(config.NewDefaultResourceResolver, component.WithName("procyonDefaultConfigResourceResolver"))
	component.Register(config.NewFileLoader, component.WithName("procyonConfigFileLoader"))
	component.Register(config.NewImporter, component.WithName("procyonConfigImporter"))
	// runtime/property
	component.Register(property.NewYamlSourceLoader, component.WithName("procyonYamlPropertySourceLoader"))
	// runtime
	component.Register(runtime.NewServerProperties, component.WithPrototypeScope())
	component.Register(runtime.NewLifecycleProperties, component.WithSingletonScope())
	return nil
}
