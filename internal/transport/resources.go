package transport

import (
	"context"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	createCh chan<- resources.Resource
	removeCh chan<- resources.Descriptor
}

func NewResources(createCh chan<- resources.Resource, removeCh chan<- resources.Descriptor) *Resources {
	return &Resources{
		createCh: createCh,
		removeCh: removeCh,
	}
}

func (r Resources) Create(_ context.Context, resource *api.Resource) (*api.Response, error) {
	convertedResource, err := protoconv.ResourceFromProto(resource)
	if err != nil {
		return nil, err
	}

	r.createCh <- convertedResource

	return &api.Response{}, nil
}

func (r Resources) Remove(_ context.Context, descriptor *api.Descriptor) (*api.Response, error) {
	kind, err := protoconv.KindFromProto(descriptor.Kind)
	if err != nil {
		return nil, err
	}

	r.removeCh <- resources.Descriptor{
		Name: descriptor.Name,
		Kind: kind,
	}

	return &api.Response{}, nil
}
