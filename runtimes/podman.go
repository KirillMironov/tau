package runtimes

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
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
	quiet := true

	_, err := images.Pull(p.ctx, container.Image, &images.PullOptions{Quiet: &quiet})
	if err != nil {
		return err
	}

	spec := specgen.NewSpecGenerator(container.Image, false)
	spec.Remove = true
	spec.Name = container.Name
	spec.Command = strings.Split(container.Command, " ")

	_, err = containers.CreateWithSpec(p.ctx, spec, nil)
	if err != nil {
		return err
	}

	return containers.Start(p.ctx, container.Name, nil)
}

func (p Podman) Remove(containerName string) error {
	err := containers.Stop(p.ctx, containerName, nil)
	if err != nil {
		return err
	}

	volumes := true

	err = containers.Remove(p.ctx, containerName, &containers.RemoveOptions{Volumes: &volumes})
	if err != nil && !strings.HasSuffix(err.Error(), define.ErrNoSuchCtr.Error()) {
		return err
	}

	return nil
}

func PodmanRootlessSocket() string {
	return "unix://" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
}
