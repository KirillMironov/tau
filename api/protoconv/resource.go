package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func ResourceFromProto(resource *api.Resource) (tau.Resource, error) {
	switch v := resource.GetKind().(type) {
	case *api.Resource_Container:
		return ContainerFromProto(v.Container), nil
	case *api.Resource_Pod:
		return PodFromProto(resource.GetPod()), nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}

func ResourceToProto(resource tau.Resource) (*api.Resource, error) {
	switch v := resource.(type) {
	case *resources.Container:
		return &api.Resource{Kind: &api.Resource_Container{Container: ContainerToProto(*v)}}, nil
	case *resources.Pod:
		return &api.Resource{Kind: &api.Resource_Pod{Pod: PodToProto(*v)}}, nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}

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

func StatusToProto(status tau.Status) (*api.Status, error) {
	state, err := StateToProto(status.State)
	if err != nil {
		return nil, err
	}

	return &api.Status{State: state}, nil
}
