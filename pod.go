package tau

import "github.com/hashicorp/go-multierror"

type Pod struct {
	Kind       string      `validate:"required,eq=pod"`
	Name       string      `validate:"required"`
	Containers []Container `validate:"required,dive"`
}

func (p Pod) Create(runtime ContainerRuntime) error {
	for _, container := range p.Containers {
		err := runtime.Start(&container)
		if err != nil {
			_ = p.Delete(runtime)
			return err
		}
	}

	return nil
}

func (p Pod) Delete(runtime ContainerRuntime) (err error) {
	for _, container := range p.Containers {
		err = multierror.Append(err, runtime.Remove(container.Id()))
	}

	return err
}

func (p Pod) Validate() error {
	return validate.Struct(p)
}
