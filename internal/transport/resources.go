package transport

import (
	"context"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
)

type Resources struct {
	createCh chan<- tau.Resource
	removeCh chan<- tau.Resource
}

func NewResources(createCh, removeCh chan<- tau.Resource) *Resources {
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

func (r Resources) Remove(_ context.Context, resource *api.Resource) (*api.Response, error) {
	convertedResource, err := protoconv.ResourceFromProto(resource)
	if err != nil {
		return nil, err
	}

	r.removeCh <- convertedResource

	return &api.Response{}, nil
}
