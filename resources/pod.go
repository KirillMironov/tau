package resources

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/KirillMironov/tau/runtimes"
)

type Pod struct {
	Name       string
	Containers []Container
	status     Status
}

func (p *Pod) Create(runtime runtimes.ContainerRuntime) error {
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

func (p *Pod) Remove(runtime runtimes.ContainerRuntime) error {
	err := p.validate()
	if err != nil {
		return err
	}

	for _, container := range p.Containers {
		err = multierror.Append(err, container.Remove(runtime))
	}

	return err
}

func (p *Pod) Descriptor() Descriptor {
	return Descriptor{
		Name: p.Name,
		Kind: KindPod,
	}
}

func (p *Pod) Status() Status {
	return p.status
}

func (p *Pod) SetState(state State) {
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

type PodAlias Pod

type PodGob struct {
	*PodAlias
	Status Status
}

func (p *Pod) MarshalBinary() ([]byte, error) {
	var (
		pod = PodGob{
			PodAlias: (*PodAlias)(p),
			Status:   p.status,
		}
		buf = new(bytes.Buffer)
	)

	err := gob.NewEncoder(buf).Encode(pod)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Pod) UnmarshalBinary(data []byte) error {
	var pod = &PodGob{
		PodAlias: (*PodAlias)(p),
	}

	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&pod)
	if err != nil {
		return err
	}

	p.status = pod.Status

	return nil
}
