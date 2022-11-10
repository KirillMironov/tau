package tau

import (
	"context"

	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/containers/podman/v2/pkg/bindings/images"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/specgen"
)

type Pod struct {
	Name  string
	Image string
}

func (p Pod) Run(ctx context.Context) error {
	_, err := images.Pull(ctx, p.Image, entities.ImagePullOptions{Quiet: true})
	if err != nil {
		return err
	}

	spec := specgen.NewSpecGenerator(p.Image, false)
	spec.Name = p.Name

	_, err = containers.CreateWithSpec(ctx, spec)
	if err != nil {
		return err
	}

	return containers.Start(ctx, p.Name, nil)
}

func (p Pod) Delete(ctx context.Context) error {
	err := containers.Stop(ctx, p.Name, nil)
	if err != nil {
		return err
	}

	volumes := true
	return containers.Remove(ctx, p.Name, nil, &volumes)
}
