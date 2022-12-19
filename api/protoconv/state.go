package protoconv

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

func StateToProto(state tau.State) (api.State, error) {
	switch state {
	case tau.StatePending:
		return api.State_STATE_PENDING, nil
	case tau.StateRunning:
		return api.State_STATE_RUNNING, nil
	case tau.StateSucceeded:
		return api.State_STATE_SUCCEEDED, nil
	case tau.StateFailed:
		return api.State_STATE_FAILED, nil
	default:
		return api.State_STATE_UNSPECIFIED, fmt.Errorf("unexpected resource state: %s", state)
	}
}
