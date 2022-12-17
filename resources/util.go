package resources

import (
	"fmt"

	"github.com/KirillMironov/tau"
)

// ResourceByKind is used for unmarshalling resources based on their kind.
func ResourceByKind(kind tau.Kind) (tau.Resource, error) {
	switch kind {
	case tau.KindContainer:
		return &Container{}, nil
	case tau.KindPod:
		return &Pod{}, nil
	default:
		return nil, fmt.Errorf("unexpected resource kind: %s", kind)
	}
}
