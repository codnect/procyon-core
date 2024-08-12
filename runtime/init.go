package runtime

import (
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/condition"
	"github.com/codnect/procyoncore/event"
)

func init() {
	// runtime
	component.Register(newDefaultCustomizer, component.Named("procyonDefaultRuntimeCustomizer"))
	component.Register(NewEventMulticaster, component.Named("procyonRuntimeEventMulticaster")).
		ConditionalOn(condition.OnMissingType[event.Multicaster]())
}
