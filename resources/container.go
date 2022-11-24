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

func (c Container) ID() string {
	return c.Name
}

func (c Container) Validate() error {
	switch {
	case c.Name == "":
		return errors.New("name is required")
	case c.Image == "":
		return errors.New("image is required")
	default:
		return nil
	}
}

func (c Container) Create(runtime runtimes.ContainerRuntime) error {
	return runtime.Start(runtimes.Container(c))
}

func (c Container) Remove(runtime runtimes.ContainerRuntime) error {
	return runtime.Remove(c.Name)
}
