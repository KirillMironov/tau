package resources

const (
	StatePending     State = "pending"
	StateRunning     State = "running"
	StateTerminating State = "terminating"
	StateSucceeded   State = "succeeded"
	StateFailed      State = "failed"
)

type State string
