package resources

import (
	"github.com/KirillMironov/tau"
	"github.com/hashicorp/go-multierror"
)

type Pod struct {
	Name       string          `validate:"required"`
	Containers []tau.Container `validate:"required,dive"`
}

func (p Pod) Create(runtime tau.ContainerRuntime) error {
	for _, container := range p.Containers {
		err := runtime.Start(&container)
		if err != nil {
			_ = p.Delete(runtime)
			return err
		}
	}

	return nil
}

func (p Pod) Delete(runtime tau.ContainerRuntime) (err error) {
	for _, container := range p.Containers {
		err = multierror.Append(err, runtime.Remove(container.Id()))
	}

	return err
}
