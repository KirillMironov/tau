package resources

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/KirillMironov/tau"
)

type containerStatus struct {
	CreatedAt int64
}

type Container struct {
	tau.Container
	state  tau.State
	status containerStatus
}

func (c *Container) Create(runtime tau.ContainerRuntime) error {
	if err := c.validate(); err != nil {
		return err
	}

	c.state = tau.StatePending
	c.status.CreatedAt = time.Now().Unix()

	return runtime.Start(tau.Container{
		Name:    c.Name,
		Image:   c.Image,
		Command: c.Command,
	})
}

func (c *Container) Remove(runtime tau.ContainerRuntime) error {
	if err := c.validate(); err != nil {
		return err
	}

	return runtime.Remove(c.Name)
}

func (c *Container) UpdateStatus(runtime tau.ContainerRuntime) error {
	if err := c.validate(); err != nil {
		return err
	}

	runtimeState, err := runtime.State(c.Name)
	if err != nil {
		return err
	}

	switch runtimeState {
	case tau.ContainerStateRunning:
		c.state = tau.StateRunning
	case tau.ContainerStateSucceeded:
		c.state = tau.StateSucceeded
	case tau.ContainerStateFailed:
		c.state = tau.StateFailed
	}

	return nil
}

func (c *Container) Descriptor() tau.Descriptor {
	return tau.Descriptor{
		Name: c.Name,
		Kind: tau.KindContainer,
	}
}

func (c *Container) State() tau.State {
	return c.state
}

func (c *Container) Status() []tau.StatusEntry {
	return []tau.StatusEntry{
		{Title: "crated at", Value: time.Unix(c.status.CreatedAt, 0).String()},
	}
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
	State     tau.State
	Status    containerStatus
}

func (c *Container) MarshalBinary() ([]byte, error) {
	container := containerGob{
		Container: (*containerAlias)(c),
		State:     c.state,
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

	c.state = container.State
	c.status = container.Status

	return nil
}
