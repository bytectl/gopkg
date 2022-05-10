package mqtt

import (
	"context"
	"time"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ Context = (*wrapper)(nil)

// Context is an MQTT Context.
type Context interface {
	context.Context
	Client() pmqtt.Client
	Message() pmqtt.Message
	Reset(pmqtt.Client, pmqtt.Message)
}

func WithContext(ctx context.Context) Context {
	return &wrapper{ctx: ctx}
}

type wrapper struct {
	ctx    context.Context
	client pmqtt.Client
	msg    pmqtt.Message
}

func (c *wrapper) Client() pmqtt.Client   { return c.client }
func (c *wrapper) Message() pmqtt.Message { return c.msg }

func (c *wrapper) Reset(client pmqtt.Client, msg pmqtt.Message) {
	c.client = client
	c.msg = msg
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
