package runtime

import "github.com/codnect/procyoncore/component"

func init() {
	// runtime
	component.Register(newDefaultCustomizer, component.Named("procyonDefaultRuntimeCustomizer"))
}
