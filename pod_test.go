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

func TestPod_Run(t *testing.T) {
	ctx, pod := setup(t)

	t.Cleanup(func() { _ = pod.Delete(ctx) })

	err := pod.Run(ctx)
	require.NoError(t, err)

	data, err := containers.Inspect(ctx, pod.Name, nil)
	require.NoError(t, err)

	assert.Equal(t, pod.Name, data.Name)
	assert.Equal(t, pod.Image, data.ImageName)
	assert.False(t, data.State.Dead)
}

func TestPod_Delete(t *testing.T) {
	ctx, pod := setup(t)

	err := pod.Run(ctx)
	require.NoError(t, err)

	exists, err := containers.Exists(ctx, pod.Name, false)
	require.NoError(t, err)
	assert.True(t, exists)

	err = pod.Delete(ctx)
	require.NoError(t, err)

	exists, err = containers.Exists(ctx, pod.Name, false)
	require.NoError(t, err)
	assert.False(t, exists)
}

func setup(t *testing.T) (context.Context, Pod) {
	t.Helper()

	socket := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"

	ctx, err := bindings.NewConnection(context.Background(), socket)
	require.NoError(t, err)

	return ctx, Pod{
		Name:  t.Name(),
		Image: image,
	}
}
