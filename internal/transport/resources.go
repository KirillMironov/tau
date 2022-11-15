package transport

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resource"
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

func (r Resources) Create(_ context.Context, request *api.Request) (*api.Response, error) {
	target, err := r.resourceByKind(request.Data)
	if err != nil {
		return nil, err
	}

	r.createCh <- target

	return &api.Response{}, nil
}

func (r Resources) Remove(_ context.Context, request *api.Request) (*api.Response, error) {
	target, err := r.resourceByKind(request.Data)
	if err != nil {
		return nil, err
	}

	r.removeCh <- target

	return &api.Response{}, nil
}

func (r Resources) resourceByKind(data []byte) (target tau.Resource, _ error) {
	var input struct {
		Kind string
	}

	err := toml.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	switch input.Kind {
	case resource.KindPod:
		target = &resource.Pod{}
	default:
		return nil, fmt.Errorf("unexpected resource kind: %s", input.Kind)
	}

	return target, toml.Unmarshal(data, &target)
}
