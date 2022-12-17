package tomlutil

import (
	"github.com/BurntSushi/toml"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/resources"
)

func UnmarshalByKind(data []byte) (tau.Resource, error) {
	var input struct {
		Kind tau.Kind
	}

	err := toml.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	resource, err := resources.ResourceByKind(input.Kind)
	if err != nil {
		return nil, err
	}

	return resource, toml.Unmarshal(data, resource)
}
