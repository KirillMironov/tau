package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

func KindFromProto(kind api.Kind) (tau.Kind, error) {
	switch kind {
	case api.Kind_KIND_CONTAINER:
		return tau.KindContainer, nil
	case api.Kind_KIND_POD:
		return tau.KindPod, nil
	default:
		return "", fmt.Errorf("unexpected resource kind: %s", kind)
	}
}
