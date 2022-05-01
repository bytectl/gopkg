package amqp

import (
	"context"
	"testing"
	"time"

	ramqp "github.com/rabbitmq/amqp091-go"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	var opts = []ServerOption{}

	opts = append(opts, Url("amqp://guest:guest@localhost:5672/"))
	opts = append(opts, OnConnectHandler(func(conn *ramqp.Connection) {
		t.Logf("amqp connected: %v", conn)

	}))
	opts = append(opts, ConnectionLostHandler(func(conn *ramqp.Connection, err *ramqp.Error) {
		t.Errorf("amqp connection lost: %v,err(%v)", conn, err)
	}))
	srv := NewServer(opts...)

	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 15)
	if srv.Stop(ctx) != nil {
		t.Errorf("expected nil got %v", srv.Stop(ctx))
	}
}
