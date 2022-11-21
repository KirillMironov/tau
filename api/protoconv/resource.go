package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func ResourceFromProto(resource *api.Resource) (tau.Resource, error) {
	switch v := resource.GetKind().(type) {
	case *api.Resource_Pod:
		return PodFromProto(resource.GetPod()), nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}

func ResourceToProto(resource tau.Resource) (*api.Resource, error) {
	switch v := resource.(type) {
	case resources.Pod:
		return &api.Resource{Kind: &api.Resource_Pod{Pod: PodToProto(v)}}, nil
	case *resources.Pod:
		return &api.Resource{Kind: &api.Resource_Pod{Pod: PodToProto(*v)}}, nil
	default:
		return nil, fmt.Errorf("unexpected resource type: %T", v)
	}
}