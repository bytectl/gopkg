// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// The registered topic, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic topic segments. They match anything until the
// next '/' or the topic end:
//  topic: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := ps.ByName("user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  thirdKey   := ps[2].Key   // the name of the 3rd parameter
//  thirdValue := ps[2].Value // the value of the 3rd parameter
package mux

import (
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Handle is a function that can be registered to a route to handle MQTT
// requests. Like MQTT.HandlerFunc, but has a third parameter for the values of
// wildcards (topic variables).
type HandlerFunc func(c mqtt.Client, msg mqtt.Message, ps *Params)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	root           *node
	paramsPool     sync.Pool
	maxParams      uint16
	NotFoundHandle HandlerFunc
}

// New returns a new initialized Router.
// topic auto-correction, including trailing slashes, is enabled by default.
func NewRouter() *Router {
	return &Router{}
}

func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

// Handle registers the handler for the given pattern.
func (r *Router) Handle(topic string, handle HandlerFunc) {
	if len(topic) < 1 {
		panic("router: topic must not be empty")
	}
	if handle == nil {
		panic("handle must not be nil")
	}
	// drop share-subscribe fields
	if strings.HasPrefix(topic, "$share/") {
		topic = strings.Join(strings.Split(topic, "/")[2:], "/")
	}
	topic = strings.TrimPrefix(topic, "$queue/")
	if topic[0] != '/' {
		// fix
		topic = "/" + topic
	}
	// add route
	varsCount := uint16(0)
	if r.root == nil {
		r.root = new(node)
	}
	r.root.addRoute(topic, handle)
	// Update maxParams
	if paramsCount := countParams(topic); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}
	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

// ServeMQTT makes the router implement the mqtt.MessageHandle interface.
func (r *Router) ServeMQTT(c mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	if topic[0] != '/' {
		// fix
		topic = "/" + topic
	}
	if r.root == nil {
		if r.NotFoundHandle != nil {
			r.NotFoundHandle(c, msg, nil)
		}
		return
	}
	handle, ps, _ := r.root.getValue(topic, r.getParams)
	if handle == nil {
		if r.NotFoundHandle != nil {
			r.NotFoundHandle(c, msg, nil)
		}
		return
	}

	if ps != nil {
		handle(c, msg, ps)
		// note: handle must before putParams
		r.putParams(ps)
	} else {
		handle(c, msg, nil)
	}
}
