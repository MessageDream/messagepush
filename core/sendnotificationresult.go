package core

import (
	"time"
)

type SendNotificationResult struct {
	Notification          *Notification
	Error                 error
	ShouldRequeue         bool
	CountsAsRequeue       bool
	OldSubscriptionId     string
	NewSubscriptionId     string
	SubscriptionExpiryUtc time.Time
	IsSubscriptionExpired bool
}

func (this *SendNotificationResult) IsSuccess() bool {
	return this.Error == nil
}
