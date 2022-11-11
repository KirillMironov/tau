package tau

import (
	"context"
	"os"
	"testing"

	"github.com/containers/podman/v2/pkg/bindings"
	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const image = "docker.io/library/traefik:1.7.34"

func TestContainer_Start(t *testing.T) {
	ctx, container := setup(t)

	t.Cleanup(func() { _ = container.Delete(ctx) })

	err := container.Start(ctx)
	require.NoError(t, err)

	data, err := containers.Inspect(ctx, container.id, nil)
	require.NoError(t, err)

	assert.Equal(t, container.id, data.ID)
	assert.Equal(t, container.Image, data.ImageName)
	assert.False(t, data.State.Dead)
}

func TestContainer_Delete(t *testing.T) {
	ctx, container := setup(t)

	err := container.Start(ctx)
	require.NoError(t, err)

	exists, err := containers.Exists(ctx, container.id, false)
	require.NoError(t, err)
	assert.True(t, exists)

	err = container.Delete(ctx)
	require.NoError(t, err)

	exists, err = containers.Exists(ctx, container.id, false)
	require.NoError(t, err)
	assert.False(t, exists)
}

func setup(t *testing.T) (context.Context, Container) {
	t.Helper()

	socket := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"

	ctx, err := bindings.NewConnection(context.Background(), socket)
	require.NoError(t, err)

	return ctx, Container{Image: image}
}
