package mqtt

import (
	"context"
	"sync"

	"github.com/bytectl/gopkg/transport/mqtt/mux"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandlerFunc defines a function to serve MQTT requests.
type HandlerFunc func(Context)

// Router is an MQTT router.
type Router struct {
	pool    sync.Pool
	srv     *Server
	filters []FilterFunc
}

func newRouter(srv *Server, filters ...FilterFunc) *Router {
	r := &Router{
		srv:     srv,
		filters: filters,
	}
	r.pool.New = func() interface{} {
		return &wrapper{router: r, ctx: context.Background()}
	}
	return r
}

// Handle registers a new route with a matcher for the Topic.
func (r *Router) Handle(topic string, h HandlerFunc, filters ...FilterFunc) {

	next := Handler(mux.HandlerFunc(func(c mqtt.Client, msg mqtt.Message) {
		ctx := r.pool.Get().(Context)
		ctx.Reset(c, msg)
		h(ctx)
		ctx.Reset(nil, nil)
		r.pool.Put(ctx)
	}))
	next = FilterChain(filters...)(next)
	next = FilterChain(r.filters...)(next)
	r.srv.router.Handle(topic, next.ServeMQTT)
}
