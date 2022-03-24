package empty

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// Registry is empty registry.
type Empty struct {
}

func New() *Empty {
	return &Empty{}
}

func (e *Empty) Register(ctx context.Context, service *registry.ServiceInstance) error {
	log.Debugf("empty  registry register service: %v", service)
	return nil
}

func (e *Empty) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	log.Debugf("empty registry  deregister service: %v", service)
	return nil
}

func (e *Empty) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return []*registry.ServiceInstance{}, nil
}

func (e *Empty) Next() ([]*registry.ServiceInstance, error) {
	return []*registry.ServiceInstance{}, nil
}
func (e *Empty) Stop() error {
	return nil
}
