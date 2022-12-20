package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
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
