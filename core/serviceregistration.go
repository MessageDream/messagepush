package core

import (
	"reflect"
)

type ServiceRegistration struct {
	Service          *PushServicer
	AppKey           string
	NotificationType reflect.Type
}

func CreateServiceRegistration(service *PushServicer, appKey string, notificationType reflect.Type) *ServiceRegistration {
	return &ServiceRegistration{
		AppKey:           appKey,
		Service:          service,
		NotificationType: notificationType,
	}
}
