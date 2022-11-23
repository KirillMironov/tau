package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func ResourceFromProto(resource *api.Resource) (resources.Resource, error) {
	switch v := resource.GetKind().(type) {
	case *api.Resource_Container:
		return ContainerFromProto(v.Container), nil
	case *api.Resource_Pod:
		return PodFromProto(resource.GetPod()), nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}

func ResourceToProto(resource resources.Resource) (*api.Resource, error) {
	switch v := resource.(type) {
	case resources.Container:
		return &api.Resource{Kind: &api.Resource_Container{Container: ContainerToProto(v)}}, nil
	case *resources.Container:
		return &api.Resource{Kind: &api.Resource_Container{Container: ContainerToProto(*v)}}, nil
	case resources.Pod:
		return &api.Resource{Kind: &api.Resource_Pod{Pod: PodToProto(v)}}, nil
	case *resources.Pod:
		return &api.Resource{Kind: &api.Resource_Pod{Pod: PodToProto(*v)}}, nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}
