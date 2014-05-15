package core

import (
	"errors"
	"reflect"
	"time"
)

type PushBroker struct {
	OnChannelCreated            func(sender interface{}, pushChannel *PushChanneler)
	OnChannelDestroyed          func(sender interface{})
	OnNotificationSent          func(sender interface{}, notification *Notification)
	OnNotificationFailed        func(sender interface{}, notification *Notification, err error)
	OnNotificationRequeue       func(sender interface{}, eventArgs *NotificationRequeueEventArgs)
	OnChannelException          func(sender interface{}, pushChannel *PushChanneler, err error)
	OnServiceException          func(sender interface{}, err error)
	OnDeviceSubscriptionExpired func(sender interface{}, expiredSubscriptionId string, expirationDateUtc time.Time, notification *Notification)
	OnDeviceSubscriptionChanged func(sender interface{}, oldSubscriptionId string, newSubscriptionId string, notification *Notification)

	serviceRegistrations []ServiceRegistration
}

func CreatPushBroker() *PushBroker {
	return &PushBroker{serviceRegistrations: make([]ServiceRegistration, 10)}
}

//Registers the service to be eligible to handle queued notifications of the specified type
func (this *PushBroker) RegisterService(pushService *PushServicer, appKey string, raiseErrorOnDuplicateRegistrations bool, notificationType reflect.Type) error {
	if raiseErrorOnDuplicateRegistrations && len(this.GetRegistrations(notificationType, appKey)) > 0 {
		return errors.New("There's already a service registered to handle " + notificationType.Name() + " notification types for the Application Id: " + appKey + "[ANY].  If you want to register the service anyway, pass in the raiseErrorOnDuplicateRegistrations=true parameter to this method.")
	}
	registration := CreateServiceRegistration(pushService, appKey, notificationType)

	this.serviceRegistrations = append(this.serviceRegistrations, *registration)

	(*pushService).OnChannelCreated(this.OnChannelCreated)
	(*pushService).OnChannelDestroyed(this.OnChannelDestroyed)
	(*pushService).OnChannelException(this.OnChannelException)
	(*pushService).OnDeviceSubscriptionExpired(this.OnDeviceSubscriptionExpired)
	(*pushService).OnNotificationFailed(this.OnNotificationFailed)
	(*pushService).OnNotificationSent(this.OnNotificationSent)
	(*pushService).OnNotificationRequeue(this.OnNotificationRequeue)
	(*pushService).OnServiceException(this.OnServiceException)
	(*pushService).OnDeviceSubscriptionChanged(this.OnDeviceSubscriptionChanged)
	return nil
}

// Queues the notification to all services registered for this type of notification with a matching appKey
func (this *PushBroker) QueueNotification(notification *Notification, notificationType reflect.Type, appKey string) error {
	var services = this.GetRegistrations(notificationType, appKey)

	if services == nil || len(services) <= 0 {
		return errors.New("There are no Registered Services that handle this type of Notification")
	}
	for _, v := range services {
		v.QueueNotification(notification)
	}
	return nil
}

// Gets all the registered services
func (this *PushBroker) GetAllRegistrations() []PushServicer {

	services := make([]PushServicer, 10)
	for _, v := range this.serviceRegistrations {
		services = append(services, *(v.Service))
	}
	return services
}

// Gets all the registered services for the given notification type and application identifier
func (this *PushBroker) GetRegistrations(notificationType reflect.Type, appKey string) []PushServicer {
	services := make([]PushServicer, 10)
	if appKey != "" && notificationType != nil {
		for _, v := range this.serviceRegistrations {
			if v.AppKey == appKey && v.NotificationType == notificationType {
				services = append(services, *(v.Service))
			}
		}
	} else if appKey != "" {
		for _, v := range this.serviceRegistrations {
			if v.AppKey == appKey {
				services = append(services, *(v.Service))
			}
		}

	} else if notificationType != nil {
		for _, v := range this.serviceRegistrations {
			if v.NotificationType == notificationType {
				services = append(services, *(v.Service))
			}
		}
	}
	return services
}

// Stops and removes all registered services for the given application identifier and notification type
func (this *PushBroker) StopAllServices(notificationType reflect.Type, appKey string, waitForQueuesToFinish bool) {

	stopping := make([]ServiceRegistration, 10)
	if notificationType == nil && appKey == "" {

		stopping = append(stopping, this.serviceRegistrations...)
		this.serviceRegistrations = append(this.serviceRegistrations[:0], this.serviceRegistrations[len(this.serviceRegistrations)-1:]...)
	} else if appKey == "" {
		for i, v := range this.serviceRegistrations {
			if v.NotificationType == notificationType {
				stopping = append(stopping, this.serviceRegistrations[i])
				this.serviceRegistrations = append(this.serviceRegistrations[:i+1], this.serviceRegistrations[i+2:]...)
			}
		}
	} else if notificationType == nil {
		for i, v := range this.serviceRegistrations {
			if v.NotificationType == notificationType {
				stopping = append(stopping, this.serviceRegistrations[i])
				this.serviceRegistrations = append(this.serviceRegistrations[:i+1], this.serviceRegistrations[i+2:]...)
			}
		}
	} else {
		for i, v := range this.serviceRegistrations {
			if v.NotificationType == notificationType && v.AppKey == appKey {
				stopping = append(stopping, this.serviceRegistrations[i])
				this.serviceRegistrations = append(this.serviceRegistrations[:i+1], this.serviceRegistrations[i+2:]...)
			}
		}
	}

	if len(stopping) != 0 {
		for _, v := range stopping {
			go StopService(v, waitForQueuesToFinish)
		}
	}
}

func StopService(sr ServiceRegistration, waitForQueuesToFinish bool) {
	(*(sr.Service)).Stop(waitForQueuesToFinish)
	//sr.Service.OnChannelCreated -= OnChannelCreated;
	//sr.Service.OnChannelDestroyed -= OnChannelDestroyed;
	//sr.Service.OnChannelException -= OnChannelException;
	//sr.Service.OnDeviceSubscriptionExpired -= OnDeviceSubscriptionExpired;
	//sr.Service.OnNotificationFailed -= OnNotificationFailed;
	//sr.Service.OnNotificationSent -= OnNotificationSent;
	//sr.Service.OnServiceException -= OnServiceException;
	//sr.Service.OnDeviceSubscriptionChanged -= OnDeviceSubscriptionChanged;
}
