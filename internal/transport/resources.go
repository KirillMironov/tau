package transport

import (
	"context"

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
	var target tau.Resource

	switch request.Kind {
	case api.Kind_POD:
		target = resource.Pod{}
	}

	err := toml.Unmarshal(request.Data, &target)
	if err != nil {
		return nil, err
	}

	r.createCh <- target

	return &api.Response{}, nil
}

func (r Resources) Remove(_ context.Context, request *api.Request) (*api.Response, error) {
	var target tau.Resource

	switch request.Kind {
	case api.Kind_POD:
		target = resource.Pod{}
	}

	err := toml.Unmarshal(request.Data, &target)
	if err != nil {
		return nil, err
	}

	r.removeCh <- target

	return &api.Response{}, nil
}
