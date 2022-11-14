package runtime

import (
	"context"
	"errors"

	"github.com/KirillMironov/tau"
	"github.com/containers/podman/v2/pkg/bindings"
	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"
)

type Podman struct {
	ctx context.Context
}

func NewPodman(socket string) (*Podman, error) {
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return nil, err
	}

	return &Podman{ctx: ctx}, nil
}

func (p Podman) Start(container *tau.Container) error {
	_, err := images.Pull(p.ctx, container.Image, entities.ImagePullOptions{Quiet: true})
	if err != nil {
		return err
	}

	spec := specgen.NewSpecGenerator(container.Image, false)

	response, err := containers.CreateWithSpec(p.ctx, spec)
	if err != nil {
		return err
	}

	containerId := response.ID

	err = containers.Start(p.ctx, containerId, nil)
	if err != nil {
		return err
	}

	container.SetId(containerId)

	return nil
}

func (p Podman) Remove(containerId string) error {
	if containerId == "" {
		return errors.New("empty container id")
	}

	err := containers.Stop(p.ctx, containerId, nil)
	if err != nil {
		return err
	}

	volumes := true
	return containers.Remove(p.ctx, containerId, nil, &volumes)
}