package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func KindFromProto(kind api.Kind) (resources.Kind, error) {
	switch kind {
	case api.Kind_KIND_CONTAINER:
		return resources.KindContainer, nil
	case api.Kind_KIND_POD:
		return resources.KindPod, nil
	default:
		return "", fmt.Errorf("unexpected resource kind: %s", kind)
	}
}

func KindToProto(kind resources.Kind) (api.Kind, error) {
	switch kind {
	case resources.KindContainer:
		return api.Kind_KIND_CONTAINER, nil
	case resources.KindPod:
		return api.Kind_KIND_POD, nil
	default:
		return api.Kind_KIND_UNSPECIFIED, fmt.Errorf("unexpected resource kind: %s", kind)
	}
}
