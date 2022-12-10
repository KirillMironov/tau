package transport

import (
	"context"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	service service
}

type service interface {
	Create(resources.Resource) error
	Remove(resources.Descriptor) error
	Get(resources.Descriptor) (resources.Status, error)
}

func NewResources(service service) *Resources {
	return &Resources{service: service}
}

func (r Resources) Create(_ context.Context, resource *api.Resource) (*api.Response, error) {
	convertedResource, err := protoconv.ResourceFromProto(resource)
	if err != nil {
		return nil, err
	}

	return &api.Response{}, r.service.Create(convertedResource)
}

func (r Resources) Get(_ context.Context, descriptor *api.Descriptor) (*api.Status, error) {
	convertedDescriptor, err := protoconv.DescriptorFromProto(descriptor)
	if err != nil {
		return nil, err
	}

	status, err := r.service.Get(convertedDescriptor)
	if err != nil {
		return nil, err
	}

	return protoconv.StatusToProto(status)
}

func (r Resources) Remove(_ context.Context, descriptor *api.Descriptor) (*api.Response, error) {
	convertedDescriptor, err := protoconv.DescriptorFromProto(descriptor)
	if err != nil {
		return nil, err
	}

	return &api.Response{}, r.service.Remove(convertedDescriptor)
}
