package event

import "time"

type ApplicationEvent interface {
	EventSource() any
	EventTime() time.Time
}
