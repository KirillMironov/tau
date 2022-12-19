package tau

import "encoding"

type Kind string

const (
	KindContainer Kind = "container"
	KindPod       Kind = "pod"
)

type State string

const (
	StatePending   State = "pending"
	StateRunning   State = "running"
	StateSucceeded State = "succeeded"
	StateFailed    State = "failed"
)

type Descriptor struct {
	Name string
	Kind Kind
}

type StatusEntry struct {
	Title string
	Value string
}

type Resource interface {
	Create(ContainerRuntime) error
	Remove(ContainerRuntime) error
	UpdateStatus(ContainerRuntime) error

	Descriptor() Descriptor
	State() State
	Status() []StatusEntry

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
