package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handler interface {
	ServeMQTT(mqtt.Client, mqtt.Message)
}

// FilterFunc is a function which receives an Handler and returns another Handler.
type FilterFunc func(Handler) Handler

// FilterChain returns a FilterFunc that specifies the chained handler for MQTT Router.
func FilterChain(filters ...FilterFunc) FilterFunc {
	return func(next Handler) Handler {
		for i := len(filters) - 1; i >= 0; i-- {
			next = filters[i](next)
		}
		return next
	}
}
