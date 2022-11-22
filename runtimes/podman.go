package runtimes

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/containers/podman/v2/libpod/define"
	"github.com/containers/podman/v2/pkg/bindings"
	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"

	"github.com/KirillMironov/tau"
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
	spec.Remove = true
	spec.Command = strings.Split(container.Command, " ")

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

	err = containers.Remove(p.ctx, containerId, nil, &volumes)
	if err != nil && !strings.HasSuffix(err.Error(), define.ErrNoSuchCtr.Error()) {
		return err
	}

	return nil
}

func PodmanRootlessSocket() string {
	return "unix://" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
}
