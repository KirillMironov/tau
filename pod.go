package tau

import (
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-multierror"
)

type Pod struct {
	Kind       string      `validate:"required,eq=pod"`
	Name       string      `validate:"required"`
	Containers []Container `validate:"required,dive"`
}

func NewPod(data []byte) (pod Pod, err error) {
	err = toml.Unmarshal(data, &pod)
	if err != nil {
		return Pod{}, err
	}

	err = validate.Struct(&pod)
	if err != nil {
		return Pod{}, err
	}

	return pod, nil
}

func (p Pod) Create(runtime ContainerRuntime) error {
	for _, container := range p.Containers {
		err := container.Start(runtime)
		if err != nil {
			_ = p.Delete(runtime)
			return err
		}
	}

	return nil
}

func (p Pod) Delete(runtime ContainerRuntime) (err error) {
	for _, container := range p.Containers {
		err = multierror.Append(err, container.Remove(runtime))
	}

	return err
}
