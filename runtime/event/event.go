package event

import "time"

type ApplicationEvent interface {
	EventSource() any

	Time() time.Time
}
