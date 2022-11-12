package runtime

import (
	"os"
	"testing"

	"github.com/KirillMironov/tau"
	"github.com/containers/podman/v2/pkg/bindings/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const image = "docker.io/library/traefik:1.7.34"

func TestPodman_Start(t *testing.T) {
	var (
		container = tau.Container{
			Image: image,
		}
		runtime = setup(t)
	)

	containerId, err := runtime.Start(container)
	require.NoError(t, err)

	t.Cleanup(func() { _ = runtime.Remove(containerId) })

	data, err := containers.Inspect(runtime.ctx, containerId, nil)
	require.NoError(t, err)

	assert.Equal(t, containerId, data.ID)
	assert.Equal(t, container.Image, data.ImageName)
	assert.False(t, data.State.Dead)
}

func TestPodman_Remove(t *testing.T) {
	var (
		container = tau.Container{
			Image: image,
		}
		runtime = setup(t)
	)

	containerId, err := runtime.Start(container)
	require.NoError(t, err)

	exists, err := containers.Exists(runtime.ctx, containerId, false)
	require.NoError(t, err)
	assert.True(t, exists)

	err = runtime.Remove(containerId)
	require.NoError(t, err)

	exists, err = containers.Exists(runtime.ctx, containerId, false)
	require.NoError(t, err)
	assert.False(t, exists)
}

func setup(t *testing.T) *Podman {
	t.Helper()

	socket := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"

	podman, err := NewPodman(socket)
	require.NoError(t, err)

	return podman
}
