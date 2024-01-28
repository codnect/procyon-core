package availability

import (
	"codnect.io/procyon-core/runtime/health"
)

type State interface {
	Status() health.Status
}
