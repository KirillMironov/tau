package resources

import "github.com/KirillMironov/tau/runtimes"

type Resource interface {
	ID() string
	Create(runtimes.ContainerRuntime) error
	Remove(runtimes.ContainerRuntime) error
}
