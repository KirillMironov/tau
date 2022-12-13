package runtimes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Docker struct {
	client *client.Client
}

func NewDocker(client *client.Client) *Docker {
	return &Docker{client: client}
}

func (d Docker) Start(container Container) error {
	if container.Name == "" {
		return errors.New("container name is required")
	}

	ctx := context.Background()

	logs, err := d.client.ImagePull(ctx, container.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, _ = io.Copy(io.Discard, logs)

	config := &containertypes.Config{
		Image: container.Image,
		Cmd:   strings.Split(container.Command, " "),
	}

	response, err := d.client.ContainerCreate(ctx, config, nil, nil, nil, container.Name)
	if err != nil {
		return err
	}

	return d.client.ContainerStart(ctx, response.ID, types.ContainerStartOptions{})
}

func (d Docker) Stop(containerName string, timeout time.Duration) error {
	err := d.client.ContainerStop(context.Background(), containerName, &timeout)
	if client.IsErrNotFound(err) {
		return ErrContainerNotFound
	}

	return err
}

func (d Docker) Remove(containerName string) error {
	options := types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}
	_ = d.client.ContainerRemove(context.Background(), containerName, options)

	return nil
}

func (d Docker) State(containerName string) (ContainerState, error) {
	status, err := d.client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return 0, ErrContainerNotFound
		}
		return 0, err
	}

	state := status.State

	switch {
	case state.Running:
		return ContainerStateRunning, nil
	case state.ExitCode == 0:
		return ContainerStateSucceeded, nil
	case state.ExitCode > 0:
		return ContainerStateFailed, nil
	default:
		return 0, fmt.Errorf("unexpected container state: %s", state.Status)
	}
}
