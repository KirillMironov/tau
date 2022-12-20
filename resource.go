package tau

import (
	"encoding"
	"fmt"
)

type Kind string

const (
	KindContainer Kind = "container"
	KindPod       Kind = "pod"
)

type State int

const (
	StateCreating State = iota
	StateRunning
	StateSucceeded
	StateFailed
)

func (s State) String() string {
	switch s {
	case StateCreating:
		return "creating"
	case StateRunning:
		return "running"
	case StateSucceeded:
		return "succeeded"
	case StateFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

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
