package protoconv

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

func DescriptorFromProto(descriptor *api.Descriptor) (tau.Descriptor, error) {
	kind, err := KindFromProto(descriptor.Kind)
	if err != nil {
		return tau.Descriptor{}, err
	}

	return tau.Descriptor{
		Name: descriptor.Name,
		Kind: kind,
	}, nil
}
