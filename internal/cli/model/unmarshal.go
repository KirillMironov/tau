package model

import (
	"github.com/BurntSushi/toml"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

type Resource interface {
	Descriptor() *api.Descriptor
	ToProto() *api.Resource
}

func UnmarshalByKind(data []byte) (Resource, error) {
	var input struct {
		Kind tau.Kind
	}

	err := toml.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	var resource Resource

	switch input.Kind {
	case tau.KindContainer:
		resource = new(Container)
	case tau.KindPod:
		resource = new(Pod)
	}

	return resource, toml.Unmarshal(data, resource)
}
