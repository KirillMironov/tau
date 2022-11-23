package runtimes

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/containers/podman/v2/libpod/define"
	"github.com/containers/podman/v2/pkg/bindings"
	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"
	"github.com/sirupsen/logrus"
)

// Disable Podman bindings logging (https://github.com/containers/podman/issues/13504).
func init() {
	logrus.SetOutput(io.Discard)
}

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

func (p Podman) Start(container Container) error {
	_, err := images.Pull(p.ctx, container.Image, entities.ImagePullOptions{Quiet: true})
	if err != nil {
		return err
	}

	spec := specgen.NewSpecGenerator(container.Image, false)
	spec.Remove = true
	spec.Name = container.Name
	spec.Command = strings.Split(container.Command, " ")

	_, err = containers.CreateWithSpec(p.ctx, spec)
	if err != nil {
		return err
	}

	return containers.Start(p.ctx, container.Name, nil)
}

func (p Podman) Remove(containerName string) error {
	if containerName == "" {
		return errors.New("empty container name")
	}

	err := containers.Stop(p.ctx, containerName, nil)
	if err != nil {
		return err
	}

	volumes := true

	err = containers.Remove(p.ctx, containerName, nil, &volumes)
	if err != nil && !strings.HasSuffix(err.Error(), define.ErrNoSuchCtr.Error()) {
		return err
	}

	return nil
}

func PodmanRootlessSocket() string {
	return "unix://" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
}
