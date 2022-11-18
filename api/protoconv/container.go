package protoconv

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

func ContainersFromProto(containers []*api.Container) []tau.Container {
	target := make([]tau.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, tau.Container{
			Image:   container.Image,
			Command: container.Command,
		})
	}

	return target
}

func ContainersToProto(containers []tau.Container) []*api.Container {
	target := make([]*api.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, &api.Container{
			Image:   container.Image,
			Command: container.Command,
		})
	}

	return target
}
