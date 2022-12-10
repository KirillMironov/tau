package resources

import (
	"encoding"

	"github.com/KirillMironov/tau/runtimes"
)

type Descriptor struct {
	Name string
	Kind Kind
}

type Status struct {
	State State
}

type Resource interface {
	Create(runtimes.ContainerRuntime) error
	Remove(runtimes.ContainerRuntime) error

	Descriptor() Descriptor
	Status() Status
	SetState(State)

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
