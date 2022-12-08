package resources

import (
	"errors"

	"github.com/KirillMironov/tau/runtimes"
)

type Container struct {
	Name    string
	Image   string
	Command string
}

func (c Container) Descriptor() Descriptor {
	return Descriptor{
		Name: c.Name,
		Kind: KindContainer,
	}
}

func (c Container) Create(runtime runtimes.ContainerRuntime) error {
	err := c.validate()
	if err != nil {
		return err
	}

	return runtime.Start(runtimes.Container(c))
}

func (c Container) Remove(runtime runtimes.ContainerRuntime) error {
	err := c.validate()
	if err != nil {
		return err
	}

	return runtime.Remove(c.Name)
}

func (c Container) validate() error {
	switch {
	case c.Name == "":
		return errors.New("name is required")
	case c.Image == "":
		return errors.New("image is required")
	default:
		return nil
	}
}
