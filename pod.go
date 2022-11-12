package tau

import (
	"context"

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

func (p Pod) Deploy(ctx context.Context) error {
	for _, container := range p.Containers {
		err := container.Start(ctx)
		if err != nil {
			_ = p.Destroy(ctx)
			return err
		}
	}

	return nil
}

func (p Pod) Destroy(ctx context.Context) (err error) {
	for _, container := range p.Containers {
		err = multierror.Append(err, container.Delete(ctx))
	}

	return err
}
