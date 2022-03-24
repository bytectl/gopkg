package empty

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// Registry is empty registry.
type Registry struct {
}

func New() *Registry {
	return &Registry{}
}
func (e *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	log.Debugf("empty  registry register service: %v", service)
	return nil
}

func (e *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	log.Debugf("empty registry  deregister service: %v", service)
	return nil
}
