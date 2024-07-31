package runtime

import (
	"github.com/codnect/procyoncore/component"
)

func init() {
	// runtime
	component.Register(newCustomizer(), component.Named("procyonRuntimeCustomizer"))
}
