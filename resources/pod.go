package resources

import (
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/KirillMironov/tau/runtimes"
)

type Pod struct {
	Name       string
	Containers []Container
}

func (p Pod) ID() string {
	return p.Name
}

func (p Pod) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	for _, container := range p.Containers {
		err := container.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Pod) Create(runtime runtimes.ContainerRuntime) error {
	for _, container := range p.Containers {
		err := container.Create(runtime)
		if err != nil {
			_ = p.Remove(runtime)
			return err
		}
	}

	return nil
}

func (p Pod) Remove(runtime runtimes.ContainerRuntime) (err error) {
	for _, container := range p.Containers {
		err = multierror.Append(err, container.Remove(runtime))
	}

	return err
}
