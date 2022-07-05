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
	pool sync.Pool
	srv  *Server
}

func newRouter(srv *Server) *Router {
	r := &Router{
		srv: srv,
	}
	r.pool.New = func() interface{} {
		return &wrapper{router: r}
	}
	return r
}

// Handle registers a new route with a matcher for the Topic.
func (r *Router) Handle(topic string, h HandlerFunc) {

	next := mux.HandlerFunc(func(c mqtt.Client, msg mqtt.Message, ps *mux.Params) {
		ctx := r.pool.Get().(Context)
		ctx.Reset(context.Background(), c, msg, ps)
		h(ctx)
		ctx.Reset(nil, nil, nil, nil)
		r.pool.Put(ctx)
	})

	r.srv.router.Handle(topic, next)
}
