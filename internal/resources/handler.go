package resources

import (
	"context"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
)

type Handler struct {
	service service
}

type service interface {
	Create(tau.Resource) error
	Remove(tau.Descriptor) error
	Status(tau.Descriptor) (tau.State, []tau.StatusEntry, error)
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h Handler) Create(_ context.Context, resource *api.Resource) (*api.Response, error) {
	convertedResource, err := protoconv.ResourceFromProto(resource)
	if err != nil {
		return nil, err
	}

	return &api.Response{}, h.service.Create(convertedResource)
}

func (h Handler) Remove(_ context.Context, descriptor *api.Descriptor) (*api.Response, error) {
	convertedDescriptor, err := protoconv.DescriptorFromProto(descriptor)
	if err != nil {
		return nil, err
	}

	return &api.Response{}, h.service.Remove(convertedDescriptor)
}

func (h Handler) Status(_ context.Context, descriptor *api.Descriptor) (*api.StatusResponse, error) {
	convertedDescriptor, err := protoconv.DescriptorFromProto(descriptor)
	if err != nil {
		return nil, err
	}

	state, status, err := h.service.Status(convertedDescriptor)
	if err != nil {
		return nil, err
	}

	protoState, err := protoconv.StateToProto(state)
	if err != nil {
		return nil, err
	}

	return &api.StatusResponse{
		State:   protoState,
		Entries: protoconv.StatusEntriesToProto(status),
	}, nil
}
