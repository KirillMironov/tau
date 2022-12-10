package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func StateToProto(state resources.State) (api.State, error) {
	switch state {
	case resources.StatePending:
		return api.State_STATE_PENDING, nil
	case resources.StateRunning:
		return api.State_STATE_RUNNING, nil
	case resources.StateTerminating:
		return api.State_STATE_TERMINATING, nil
	case resources.StateSucceeded:
		return api.State_STATE_SUCCEEDED, nil
	case resources.StateFailed:
		return api.State_STATE_FAILED, nil
	default:
		return api.State_STATE_UNSPECIFIED, fmt.Errorf("unexpected resource state: %s", state)
	}
}
