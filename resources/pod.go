package resources

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/KirillMironov/tau"
)

type Pod struct {
	Name       string
	Containers []Container
	status     tau.Status
}

func (p *Pod) Create(runtime tau.ContainerRuntime) error {
	err := p.validate()
	if err != nil {
		return err
	}

	for _, container := range p.Containers {
		err = container.Create(runtime)
		if err != nil {
			_ = p.Remove(runtime)
			return err
		}
	}

	return nil
}

func (p *Pod) Remove(runtime tau.ContainerRuntime) error {
	err := p.validate()
	if err != nil {
		return err
	}

	for _, container := range p.Containers {
		err = multierror.Append(err, container.Remove(runtime))
	}

	return err
}

func (p *Pod) Descriptor() tau.Descriptor {
	return tau.Descriptor{
		Name: p.Name,
		Kind: tau.KindPod,
	}
}

func (p *Pod) Status() tau.Status {
	return p.status
}

func (p *Pod) SetState(state tau.State) {
	p.status.State = state
}

func (p *Pod) validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	for _, container := range p.Containers {
		err := container.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// podAlias is used to avoid infinite recursion during gob encoding/decoding.
type podAlias Pod

// podGob represents a gob-serializable version of Pod.
type podGob struct {
	Pod    *podAlias
	Status tau.Status
}

func (p *Pod) MarshalBinary() ([]byte, error) {
	pod := podGob{
		Pod:    (*podAlias)(p),
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

	p.status = pod.Status

	return nil
}
