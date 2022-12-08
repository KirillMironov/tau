package resources

import "github.com/KirillMironov/tau/runtimes"

type Descriptor struct {
	Name string
	Kind Kind
}

type Resource interface {
	Descriptor() Descriptor
	Create(runtimes.ContainerRuntime) error
	Remove(runtimes.ContainerRuntime) error
}
