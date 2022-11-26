package runtimes

import (
	"context"
	"io"
	"strings"

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
	ctx := context.Background()

	pullLogs, err := d.client.ImagePull(ctx, container.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer pullLogs.Close()

	_, _ = io.Copy(io.Discard, pullLogs)

	config := &containertypes.Config{
		Image: container.Image,
		Cmd:   strings.Split(container.Command, " "),
	}

	hostConfig := &containertypes.HostConfig{AutoRemove: true}

	response, err := d.client.ContainerCreate(ctx, config, hostConfig, nil, nil, container.Name)
	if err != nil {
		return err
	}

	return d.client.ContainerStart(ctx, response.ID, types.ContainerStartOptions{})
}

func (d Docker) Remove(containerName string) error {
	return d.client.ContainerStop(context.Background(), containerName, nil)
}
