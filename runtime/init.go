package runtime

import "codnect.io/procyon-core/component"

func init() {
	component.Register(newStartupListener)
}
