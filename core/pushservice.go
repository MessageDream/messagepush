package core

import (
	"errors"
	"github.com/golang/glog"
	"time"
)

type PushServicer interface {
	OnChannelCreated(handler func(sender interface{}, pushChannel *PushChanneler))
	OnChannelDestroyed(handler func(sender interface{}))
	OnNotificationSent(handler func(sender interface{}, notification *Notification))
	OnNotificationFailed(handler func(sender interface{}, notification *Notification, err error))
	OnNotificationRequeue(handler func(sender interface{}, eventArgs *NotificationRequeueEventArgs))
	OnChannelException(handler func(sender interface{}, pushChannel *PushChanneler, err error))
	OnServiceException(handler func(sender interface{}, err error))
	OnDeviceSubscriptionExpired(handler func(sender interface{}, expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification))
	OnDeviceSubscriptionChanged(handler func(sender interface{}, oldSubscriptionId string, newSubscriptionId string, notification *Notification))

	IsStopping() bool
	QueueNotification(notification *Notification) error
	Stop(waitForQueueToFinish bool)
}

type PushChanneler interface {
	SendNotification(notification *Notification, callback func(sender interface{}, result *SendNotificationResult))
}

type PushChannelFactorier interface {
	CreateChannel(channelSettings *PushChannelSettinger) *PushChanneler
}

type NotificationRequeueEventArgs struct {
	Cancel       bool
	notification *Notification
	requeueCause error
}

func (this *NotificationRequeueEventArgs) Notification() *Notification {
	return this.notification
}

func (this *NotificationRequeueEventArgs) RequeueCause() error {
	return this.requeueCause
}

type PushChannelSettinger struct {
}

type PushServiceBase struct {
	onChannelCreated            []func(sender interface{}, pushChannel *PushChanneler)
	onChannelDestroyed          []func(sender interface{})
	onNotificationSent          []func(sender interface{}, notification *Notification)
	onNotificationFailed        []func(sender interface{}, notification *Notification, err error)
	onNotificationRequeue       []func(sender interface{}, eventArgs *NotificationRequeueEventArgs)
	onChannelException          []func(sender interface{}, pushChannel *PushChanneler, err error)
	onServiceException          []func(sender interface{}, err error)
	onDeviceSubscriptionExpired []func(sender interface{}, expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification)
	onDeviceSubscriptionChanged []func(sender interface{}, oldSubscriptionId string, newSubscriptionId string, notification *Notification)

	pushChannelFactory *PushChannelFactorier
	serviceSettings    *PushServiceSettings
	channelSettings    *PushChannelSettinger
	notificationQueue  chan *Notification
	stopping           bool
}

func CreatPushService(pushChannelFactory *PushChannelFactorier, channelSettings *PushChannelSettinger, serviceSettings *PushServiceSettings) PushServicer {
	return PushServiceBase{
		pushChannelFactory: pushChannelFactory,
		serviceSettings:    serviceSettings,
		channelSettings:    channelSettings,
		notificationQueue:  make(chan *Notification, 0),
	}
}

func (this *PushServiceBase) PushChannelFactory() *PushChannelFactorier {
	return this.pushChannelFactory
}

func (this *PushServiceBase) ServiceSettings() *PushServiceSettings {
	return this.serviceSettings
}

func (this *PushServiceBase) ChannelSettings() *PushChannelSettinger {
	return this.channelSettings
}

func (this *PushServiceBase) RaiseSubscriptionExpired(expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification) {
	evt := this.onDeviceSubscriptionExpired
	if evt != nil {
		for _, v := range evt {
			v(this, expiredSubscriptionId, expirationDateUtc, notification)
		}
	}
}

func (this *PushServiceBase) RaiseServiceException(err error) {
	evt := this.onServiceException
	if evt != nil {
		for _, v := range evt {
			v(this, err)
		}
	}
}

func (this *PushServiceBase) BlockOnMessageResult() bool {
	return true
}

func (this PushServiceBase) OnChannelCreated(handler func(sender interface{}, pushChannel *PushChanneler)) {
	if this.onChannelCreated == nil {
		this.onChannelCreated = make([]func(sender interface{}, pushChannel *PushChanneler), 2)
	}
	this.onChannelCreated = append(this.onChannelCreated, handler)
}

func (this PushServiceBase) OnChannelDestroyed(handler func(sender interface{})) {
	if this.onChannelDestroyed == nil {
		this.onChannelDestroyed = make([]func(sender interface{}), 2)
	}
	this.onChannelDestroyed = append(this.onChannelDestroyed, handler)
}

func (this PushServiceBase) OnNotificationSent(handler func(sender interface{}, notification *Notification)) {
	if this.onNotificationSent == nil {
		this.onNotificationSent = make([]func(sender interface{}, notification *Notification), 2)
	}
	this.onNotificationSent = append(this.onNotificationSent, handler)
}

func (this PushServiceBase) OnNotificationFailed(handler func(sender interface{}, notification *Notification, err error)) {
	if this.onNotificationFailed == nil {
		this.onNotificationFailed = make([]func(sender interface{}, notification *Notification, err error), 2)
	}
	this.onNotificationFailed = append(this.onNotificationFailed, handler)
}

func (this PushServiceBase) OnNotificationRequeue(handler func(sender interface{}, eventArgs *NotificationRequeueEventArgs)) {
	if this.onNotificationRequeue == nil {
		this.onNotificationRequeue = make([]func(sender interface{}, eventArgs *NotificationRequeueEventArgs), 2)
	}
	this.onNotificationRequeue = append(this.onNotificationRequeue, handler)
}

func (this PushServiceBase) OnChannelException(handler func(sender interface{}, pushChannel *PushChanneler, err error)) {
	if this.onChannelException == nil {
		this.onChannelException = make([]func(sender interface{}, pushChannel *PushChanneler, err error), 2)
	}
	this.onChannelException = append(this.onChannelException, handler)
}

func (this PushServiceBase) OnServiceException(handler func(sender interface{}, err error)) {
	if this.onServiceException == nil {
		this.onServiceException = make([]func(sender interface{}, err error), 2)
	}
	this.onServiceException = append(this.onServiceException, handler)
}

func (this PushServiceBase) OnDeviceSubscriptionExpired(handler func(sender interface{}, expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification)) {
	if this.onDeviceSubscriptionExpired == nil {
		this.onDeviceSubscriptionExpired = make([]func(sender interface{}, expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification), 2)
	}
	this.onDeviceSubscriptionExpired = append(this.onDeviceSubscriptionExpired, handler)
}

func (this PushServiceBase) OnDeviceSubscriptionChanged(handler func(sender interface{}, oldSubscriptionId string, newSubscriptionId string, notification *Notification)) {
	if this.onDeviceSubscriptionChanged == nil {
		this.onDeviceSubscriptionChanged = make([]func(sender interface{}, oldSubscriptionId string, newSubscriptionId string, notification *Notification), 2)
	}
	this.onDeviceSubscriptionChanged = append(this.onDeviceSubscriptionChanged, handler)
}

func (this PushServiceBase) IsStopping() bool {
	return this.stopping
}

func (this PushServiceBase) Stop(waitForQueueToFinish bool) {
	this.stopping = true
	//started := DateTime.UtcNow;

	//		if waitForQueueToFinish{
	//			Log.Info ("Waiting for Queue to Finish");

	//			while (len(this.queuedNotifications) > 0)
	//				Thread.Sleep(100);

	//			Log.Info("Queue Emptied.");
	//		}

	//		Log.Info ("Stopping CheckScale Timer");

	//		//Stop the timer for checking scale
	//		if (this.timerCheckScale != null)
	//			this.timerCheckScale.Change(Timeout.Infinite, Timeout.Infinite);

	//		Log.Info ("Stopping all Channels");

	//		//Stop all channels
	//		Parallel.ForEach(channels, c => c.Dispose());

	//		this.channels.Clear();

	//		this.cancelTokenSource.Cancel();

	//		Log.Info("PushServiceBase->DISPOSE.");
}

func (this PushServiceBase) QueueNotification(notification *Notification) error {
	return this.queueNotification(notification, false, false)
}

func (this *PushServiceBase) queueNotification(notification *Notification, countsAsRequeue, ignoreStoppingChannel bool) error {
	notification.EnqueuedTimestamp = time.Now().UTC()
	if this.serviceSettings.MaxNotificationRequeues < 0 || notification.QueuedCount <= this.serviceSettings.MaxNotificationRequeues {
		notification.EnqueuedTimestamp = time.Now().UTC()

		this.notificationQueue <- notification

		if countsAsRequeue {
			notification.QueuedCount++
		}
	} else {
		evt := this.onNotificationFailed
		if evt != nil {
			for _, v := range evt {
				v(this, notification, errors.New("The maximum number of Send attempts to send the notification was reached!"))
			}
		}

		glog.Info("Notification ReQueued Too Many Times:%d", notification.QueuedCount)
	}

	return nil

}
