package resources

import "github.com/KirillMironov/tau/runtimes"

type Resource interface {
	ID() string
	Validate() error
	Create(runtimes.Runtime) error
	Delete(runtimes.Runtime) error
}
