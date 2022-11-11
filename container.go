package tau

import (
	"context"
	"errors"

	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"
)

type Container struct {
	Image   string `validate:"required"`
	Command []string
	id      string
}

func (c *Container) Start(ctx context.Context) error {
	_, err := images.Pull(ctx, c.Image, entities.ImagePullOptions{Quiet: true})
	if err != nil {
		return err
	}

	spec := specgen.NewSpecGenerator(c.Image, false)

	response, err := containers.CreateWithSpec(ctx, spec)
	if err != nil {
		return err
	}

	c.id = response.ID

	return containers.Start(ctx, c.id, nil)
}

func (c *Container) Delete(ctx context.Context) error {
	if c.id == "" {
		return errors.New("container not started")
	}

	err := containers.Stop(ctx, c.id, nil)
	if err != nil {
		return err
	}

	volumes := true
	return containers.Remove(ctx, c.id, nil, &volumes)
}
