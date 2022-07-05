package mqtt

import (
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/bytectl/gopkg/transport/mqtt/mux"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

var _ Context = (*wrapper)(nil)

// Context is an MQTT Context.
type Context interface {
	context.Context
	Client() pmqtt.Client
	Message() pmqtt.Message
	Reset(context.Context, pmqtt.Client, pmqtt.Message, *mux.Params)
	Middleware(middleware.Handler) middleware.Handler
	Bind(v interface{}) error
	BindVars(v interface{}) error
	Encode(v interface{}) ([]byte, error)
	EncodeErr(err error) []byte
}

type wrapper struct {
	router *Router
	ctx    context.Context
	client pmqtt.Client
	msg    pmqtt.Message
	ps     *mux.Params
}

func (c *wrapper) Client() pmqtt.Client   { return c.client }
func (c *wrapper) Message() pmqtt.Message { return c.msg }
func (c *wrapper) Middleware(h middleware.Handler) middleware.Handler {
	return middleware.Chain(c.router.srv.ms...)(h)
}
func (c *wrapper) Reset(ctx context.Context, client pmqtt.Client, msg pmqtt.Message, ps *mux.Params) {
	c.ctx = ctx
	c.client = client
	c.msg = msg
	c.ps = ps
}

func (c *wrapper) Deadline() (time.Time, bool) {
	if c.ctx == nil {
		return time.Time{}, false
	}
	return c.ctx.Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Done()
}

func (c *wrapper) Err() error {
	if c.ctx == nil {
		return context.Canceled
	}
	return c.ctx.Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Value(key)
}

type paramsKey struct{}

var pKey = paramsKey{}

func (c *wrapper) Bind(v interface{}) error { return c.router.srv.dec(c.Message().Payload(), v) }

func (c *wrapper) BindVars(v interface{}) error {
	varValues := make(url.Values)
	for _, p := range *c.ps {
		if p.Key == "" {
			continue
		}
		varValues.Add(p.Key, p.Value)
	}
	return binding.BindQuery(varValues, v)
}

func (c *wrapper) Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := c.router.srv.enc(&buf, v)
	return buf.Bytes(), err
}
func (c *wrapper) EncodeErr(err error) []byte {
	var buf bytes.Buffer
	c.router.srv.ene(&buf, err)
	return buf.Bytes()
}
