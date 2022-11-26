//go:build integration

package runtimes

import (
	"context"
	"strings"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocker_Start(t *testing.T) {
	t.Parallel()

	container, docker := dockerSetup(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := docker.Start(container)
	require.NoError(t, err)

	t.Cleanup(func() { _ = docker.Remove(container.Name) })

	data, err := docker.client.ContainerInspect(ctx, container.Name)
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(data.Name, container.Name))
	assert.Equal(t, container.Image, data.Config.Image)
	assert.False(t, data.State.Dead)
}

func TestDocker_Remove(t *testing.T) {
	t.Parallel()

	container, docker := dockerSetup(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := docker.Start(container)
	require.NoError(t, err)

	data, err := docker.client.ContainerInspect(ctx, container.Name)
	require.NoError(t, err)
	require.True(t, data.State.Running)

	err = docker.Remove(container.Name)
	require.NoError(t, err)

	data, err = docker.client.ContainerInspect(ctx, container.Name)
	require.NoError(t, err)
	assert.False(t, data.State.Running)
}

func dockerSetup(t *testing.T) (Container, *Docker) {
	t.Helper()

	container := Container{
		Name:    t.Name(),
		Image:   "docker.io/library/busybox:1.35.0",
		Command: "sleep 2m",
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)

	return container, NewDocker(cli)
}
