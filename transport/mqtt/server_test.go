package mqtt

import (
	"context"
	"testing"
	"time"
)

type testKey struct{}

type testData struct {
	Path string `json:"path"`
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(Broker("tcp://10.10.18.240:11183"))

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
