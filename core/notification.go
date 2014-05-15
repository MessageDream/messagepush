package core

import (
	"time"
)

type Notification struct {
	Tag               interface{}
	QueuedCount       int
	EnqueuedTimestamp time.Time
}

func (this *Notification) IsValidDeviceRegistrationId() bool {
	return true
}
