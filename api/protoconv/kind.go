package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func KindToProto(kind resources.Kind) (api.Kind, error) {
	switch kind {
	case resources.KindContainer:
		return api.Kind_KIND_CONTAINER, nil
	case resources.KindPod:
		return api.Kind_KIND_POD, nil
	default:
		return 0, fmt.Errorf("unexpected resource kind: %s", kind)
	}
}
