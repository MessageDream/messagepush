package core

import (
	"time"
)

type PushServiceSettings struct {
	AutoScaleChannels         bool
	MaxAutoScaleChannels      int
	MinAvgTimeToScaleChannels int64
	Channels                  int
	MaxNotificationRequeues   int
	NotificationSendTimeout   int
	IdleTimeout               time.Duration
}

func CreatDefaultPushServiceSettings() *PushServiceSettings {
	return &PushServiceSettings{
		AutoScaleChannels:         true,
		MaxAutoScaleChannels:      20,
		MinAvgTimeToScaleChannels: 100,
		Channels:                  1,
		MaxNotificationRequeues:   5,
		NotificationSendTimeout:   15000,
		IdleTimeout:               time.Duration(5),
	}
}
