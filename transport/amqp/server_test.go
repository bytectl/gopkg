package amqp

import (
	"context"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(Url("amqp://guest:guest@localhost:5672/"))

	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	if srv.Stop(ctx) != nil {
		t.Errorf("expected nil got %v", srv.Stop(ctx))
	}
}
