package resources

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/KirillMironov/tau"
)

type Container struct {
	Name    string
	Image   string
	Command string
	status  tau.Status
}

func (c *Container) Create(runtime tau.ContainerRuntime) error {
	err := c.validate()
	if err != nil {
		return err
	}

	return runtime.Start(tau.Container{
		Name:    c.Name,
		Image:   c.Image,
		Command: c.Command,
	})
}

func (c *Container) Remove(runtime tau.ContainerRuntime) error {
	err := c.validate()
	if err != nil {
		return err
	}

	return runtime.Remove(c.Name)
}

func (c *Container) Descriptor() tau.Descriptor {
	return tau.Descriptor{
		Name: c.Name,
		Kind: tau.KindContainer,
	}
}

func (c *Container) Status() tau.Status {
	return c.status
}

func (c *Container) SetState(state tau.State) {
	c.status.State = state
}

func (c *Container) validate() error {
	switch {
	case c.Name == "":
		return errors.New("name is required")
	case c.Image == "":
		return errors.New("image is required")
	default:
		return nil
	}
}

// containerAlias is used to avoid infinite recursion during gob encoding/decoding.
type containerAlias Container

// containerGob represents a gob-serializable version of Container.
type containerGob struct {
	Container *containerAlias
	Status    tau.Status
}

func (c *Container) MarshalBinary() ([]byte, error) {
	container := containerGob{
		Container: (*containerAlias)(c),
		Status:    c.status,
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(container); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Container) UnmarshalBinary(data []byte) error {
	container := containerGob{
		Container: (*containerAlias)(c),
	}

	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&container); err != nil {
		return err
	}

	c.status = container.Status

	return nil
}
