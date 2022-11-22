package tomlutil

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/resources"
)

func UnmarshalByKind(data []byte) (tau.Resource, error) {
	var resource struct {
		Kind string
	}

	err := toml.Unmarshal(data, &resource)
	if err != nil {
		return nil, err
	}

	switch resource.Kind {
	case resources.KindPod:
		pod := &resources.Pod{}
		return pod, toml.Unmarshal(data, &pod)
	default:
		return nil, fmt.Errorf("unknown resource kind: %s", resource.Kind)
	}
}
