//go:build integration

package runtimes

import (
	"testing"

	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPodman_Start(t *testing.T) {
	container, podman := podmanSetup(t)

	err := podman.Start(container)
	require.NoError(t, err)

	t.Cleanup(func() { _ = podman.Remove(container.Name) })

	data, err := containers.Inspect(podman.ctx, container.Name, nil)
	require.NoError(t, err)

	assert.Equal(t, container.Name, data.Name)
	assert.Equal(t, container.Image, data.ImageName)
	assert.False(t, data.State.Dead)
}

func TestPodman_Remove(t *testing.T) {
	container, podman := podmanSetup(t)

	err := podman.Start(container)
	require.NoError(t, err)

	exists, err := containers.Exists(podman.ctx, container.Name, nil)
	require.NoError(t, err)
	require.True(t, exists)

	err = podman.Remove(container.Name)
	require.NoError(t, err)

	exists, err = containers.Exists(podman.ctx, container.Name, nil)
	require.NoError(t, err)
	assert.False(t, exists)
}

func podmanSetup(t *testing.T) (Container, *Podman) {
	t.Helper()

	container := Container{
		Name:    t.Name(),
		Image:   "docker.io/library/busybox:1.35.0",
		Command: "sleep 2m",
	}

	podman, err := NewPodman(PodmanRootlessSocket())
	require.NoError(t, err)

	return container, podman
}
