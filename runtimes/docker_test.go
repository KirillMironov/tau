//go:build integration

package runtimes

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/KirillMironov/tau"
)

const (
	busybox = "docker.io/library/busybox:1.35.0"
	traefik = "docker.io/library/traefik:1.7.34-alpine"
)

func TestDocker_Start(t *testing.T) {
	t.Parallel()

	dockerClient := newDockerClient(t)

	runtime := NewDocker(dockerClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name      string
		container tau.Container
		wantErr   bool
	}{
		{
			name: "success",
			container: tau.Container{
				Name:    t.Name(),
				Image:   busybox,
				Command: "sleep 2m",
			},
			wantErr: false,
		},
		{
			name: "without command",
			container: tau.Container{
				Name:    t.Name(),
				Image:   traefik,
				Command: "",
			},
			wantErr: false,
		},
		{
			name: "without name",
			container: tau.Container{
				Name:    "",
				Image:   busybox,
				Command: "",
			},
			wantErr: true,
		},
		{
			name: "without image",
			container: tau.Container{
				Name:    t.Name(),
				Image:   "",
				Command: "",
			},
			wantErr: true,
		},
		{
			name: "invalid image",
			container: tau.Container{
				Name:    t.Name(),
				Image:   "invalid-image",
				Command: "",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { removeContainer(t, dockerClient, tc.container.Name) })

			err := runtime.Start(tc.container)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Errorf("Start(%#v) returned err = %v, want error presence = %v", tc.container, err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}

			status, err := dockerClient.ContainerInspect(ctx, tc.container.Name)
			if err != nil {
				t.Fatalf("failed to inspect container: %v", err)
			}

			if !status.State.Running {
				t.Error("container state is not running")
			}

			if status.Name != "/"+tc.container.Name {
				t.Errorf("container name is %q, want %q", status.Name, tc.container.Name)
			}

			if status.Config.Image != tc.container.Image {
				t.Errorf("container image is %q, want %q", status.Config.Image, tc.container.Image)
			}
		})
	}
}

func TestDocker_Stop(t *testing.T) {
	t.Parallel()

	dockerClient := newDockerClient(t)

	runtime := NewDocker(dockerClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name           string
		container      tau.Container
		startContainer bool
		wantErr        error
	}{
		{
			name: "success",
			container: tau.Container{
				Name:    t.Name(),
				Image:   traefik,
				Command: "",
			},
			startContainer: true,
			wantErr:        nil,
		},
		{
			name: "container not found",
			container: tau.Container{
				Name:    "not-found",
				Image:   traefik,
				Command: "",
			},
			startContainer: false,
			wantErr:        ErrContainerNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.startContainer {
				t.Cleanup(func() { removeContainer(t, dockerClient, tc.container.Name) })

				mustStartContainer(t, dockerClient, tc.container)
			}

			gotErr := runtime.Stop(tc.container.Name, 0)
			if gotErr != tc.wantErr {
				t.Errorf("Stop(%q, 0) returned err = %v, want error = %v", tc.container.Name, gotErr, tc.wantErr)
			}

			if tc.wantErr != nil {
				return
			}

			status, err := dockerClient.ContainerInspect(ctx, tc.container.Name)
			if err != nil {
				t.Fatalf("failed to inspect container: %v", err)
			}

			if status.State.Running {
				t.Error("container state is running, want stopped")
			}
		})
	}
}

func TestDocker_Remove(t *testing.T) {
	t.Parallel()

	dockerClient := newDockerClient(t)

	runtime := NewDocker(dockerClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name           string
		container      tau.Container
		startContainer bool
		wantErr        error
	}{
		{
			name: "success",
			container: tau.Container{
				Name:    t.Name(),
				Image:   busybox,
				Command: "sleep 2m",
			},
			startContainer: true,
			wantErr:        nil,
		},
		{
			name: "container not found",
			container: tau.Container{
				Name:    t.Name(),
				Image:   busybox,
				Command: "sleep 2m",
			},
			startContainer: false,
			wantErr:        nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.startContainer {
				t.Cleanup(func() { removeContainer(t, dockerClient, tc.container.Name) })

				mustStartContainer(t, dockerClient, tc.container)
			}

			gotErr := runtime.Remove(tc.container.Name)
			if gotErr != tc.wantErr {
				t.Errorf("Remove(%q) returned err = %v, want error = %v", tc.container.Name, gotErr, tc.wantErr)
			}

			if tc.wantErr != nil {
				return
			}

			_, err := dockerClient.ContainerInspect(ctx, tc.container.Name)
			if !client.IsErrNotFound(err) {
				t.Errorf("container %q is still present", tc.container.Name)
			}
		})
	}
}

func TestDocker_State(t *testing.T) {
	t.Parallel()

	dockerClient := newDockerClient(t)

	runtime := NewDocker(dockerClient)

	tests := []struct {
		name           string
		container      tau.Container
		startContainer bool
		wantState      tau.ContainerState
		wantErr        bool
	}{
		{
			name: "state running",
			container: tau.Container{
				Name:    t.Name(),
				Image:   busybox,
				Command: "sleep 2m",
			},
			startContainer: true,
			wantState:      tau.ContainerStateRunning,
			wantErr:        false,
		},
		{
			name: "container not found",
			container: tau.Container{
				Name:    t.Name(),
				Image:   busybox,
				Command: "sleep 2m",
			},
			startContainer: false,
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.startContainer {
				t.Cleanup(func() { removeContainer(t, dockerClient, tc.container.Name) })

				mustStartContainer(t, dockerClient, tc.container)
			}

			gotState, err := runtime.State(tc.container.Name)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Errorf("State(%q) returned err = %v, want error presence = %v", tc.container.Name, err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}

			if gotState != tc.wantState {
				t.Errorf("state = %v, want state = %v", gotState, tc.wantState)
			}
		})
	}
}

func newDockerClient(t *testing.T) *client.Client {
	t.Helper()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dockerClient.Ping(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = dockerClient.Close() })

	return dockerClient
}

func mustStartContainer(t *testing.T, dockerClient *client.Client, container tau.Container) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logs, err := dockerClient.ImagePull(ctx, container.Image, types.ImagePullOptions{})
	if err != nil {
		t.Fatalf("failed to pull image %q: %v", container.Image, err)
	}
	defer logs.Close()

	_, _ = io.Copy(io.Discard, logs)

	config := &containertypes.Config{
		Image: container.Image,
		Cmd:   strings.Split(container.Command, " "),
	}

	hostConfig := &containertypes.HostConfig{AutoRemove: true}

	_, err = dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, container.Name)
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}

	err = dockerClient.ContainerStart(ctx, container.Name, types.ContainerStartOptions{})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
}

func removeContainer(t *testing.T, dockerClient *client.Client, containerName string) {
	t.Helper()

	options := types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}
	_ = dockerClient.ContainerRemove(context.Background(), containerName, options)
}
