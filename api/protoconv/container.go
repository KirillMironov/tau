package protoconv

import (
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func ContainerFromProto(container *api.Container) *resources.Container {
	return &resources.Container{
		Name:    container.Name,
		Image:   container.Image,
		Command: container.Command,
	}
}

func ContainerToProto(container resources.Container) *api.Container {
	return &api.Container{
		Name:    container.Name,
		Image:   container.Image,
		Command: container.Command,
	}
}

func ContainersFromProto(containers []*api.Container) []resources.Container {
	target := make([]resources.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, *ContainerFromProto(container))
	}

	return target
}

func ContainersToProto(containers []resources.Container) []*api.Container {
	target := make([]*api.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, ContainerToProto(container))
	}

	return target
}
