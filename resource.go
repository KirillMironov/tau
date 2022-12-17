package tau

import "encoding"

type State string

const (
	StatePending     State = "pending"
	StateRunning     State = "running"
	StateTerminating State = "terminating"
	StateSucceeded   State = "succeeded"
	StateFailed      State = "failed"
)

type Kind string

const (
	KindContainer Kind = "container"
	KindPod       Kind = "pod"
)

type Descriptor struct {
	Name string
	Kind Kind
}

type Status struct {
	State State
}

type Resource interface {
	Create(ContainerRuntime) error
	Remove(ContainerRuntime) error

	Descriptor() Descriptor
	Status() Status
	SetState(State)

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
