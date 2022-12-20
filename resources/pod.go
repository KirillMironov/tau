package resources

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/KirillMironov/tau"
)

type podStatus struct {
	CreatedAt int64
}

type Pod struct {
	Name       string
	Containers []tau.Container
	state      tau.State
	status     podStatus
}

func (p *Pod) Create(runtime tau.ContainerRuntime) error {
	if err := p.validate(); err != nil {
		return err
	}

	p.status.CreatedAt = time.Now().Unix()

	for _, container := range p.Containers {
		if err := runtime.Start(container); err != nil {
			_ = p.Remove(runtime)
			return err
		}
	}

	return nil
}

func (p *Pod) Remove(runtime tau.ContainerRuntime) error {
	if err := p.validate(); err != nil {
		return err
	}

	var err error

	for _, container := range p.Containers {
		err = multierror.Append(err, runtime.Remove(container.Name))
	}

	return err
}

func (p *Pod) UpdateStatus(runtime tau.ContainerRuntime) error {
	if err := p.validate(); err != nil {
		return err
	}

	var hasRunning, hasSucceeded bool

	for _, container := range p.Containers {
		state, err := runtime.State(container.Name)
		if err != nil {
			return err
		}

		switch state {
		case tau.ContainerStateRunning:
			hasRunning = true
		case tau.ContainerStateSucceeded:
			hasSucceeded = true
		case tau.ContainerStateFailed:
			p.state = tau.StateFailed
			return nil
		}
	}

	p.state = tau.StateRunning

	if hasSucceeded && !hasRunning {
		p.state = tau.StateSucceeded
	}

	return nil
}

func (p *Pod) Descriptor() tau.Descriptor {
	return tau.Descriptor{
		Name: p.Name,
		Kind: tau.KindPod,
	}
}

func (p *Pod) State() tau.State {
	return p.state
}

func (p *Pod) Status() []tau.StatusEntry {
	return []tau.StatusEntry{
		{Title: "created at", Value: time.Unix(p.status.CreatedAt, 0).String()},
	}
}

func (p *Pod) validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	for _, container := range p.Containers {
		switch {
		case container.Name == "":
			return errors.New("name is required")
		case container.Image == "":
			return errors.New("image is required")
		default:
			return nil
		}
	}

	return nil
}

// podAlias is used to avoid infinite recursion during gob encoding/decoding.
type podAlias Pod

// podGob represents a gob-serializable version of Pod.
type podGob struct {
	Pod    *podAlias
	State  tau.State
	Status podStatus
}

func (p *Pod) MarshalBinary() ([]byte, error) {
	pod := podGob{
		Pod:    (*podAlias)(p),
		State:  p.state,
		Status: p.status,
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(pod); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Pod) UnmarshalBinary(data []byte) error {
	pod := podGob{
		Pod: (*podAlias)(p),
	}

	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&pod); err != nil {
		return err
	}

	p.state = pod.State
	p.status = pod.Status

	return nil
}
