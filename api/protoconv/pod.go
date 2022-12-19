package protoconv

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func PodFromProto(pod *api.Pod) *resources.Pod {
	return &resources.Pod{
		Name:       pod.Name,
		Containers: containersFromProto(pod.Containers),
	}
}

func PodToProto(pod resources.Pod) *api.Pod {
	return &api.Pod{
		Name:       pod.Name,
		Containers: containersToProto(pod.Containers),
	}
}

func containersFromProto(containers []*api.Container) []tau.Container {
	target := make([]tau.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, tau.Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
		})
	}

	return target
}

func containersToProto(containers []tau.Container) []*api.Container {
	target := make([]*api.Container, 0, len(containers))

	for _, container := range containers {
		target = append(target, &api.Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
		})
	}

	return target
}
