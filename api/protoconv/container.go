package protoconv

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func ContainerFromProto(container *api.Container) *resources.Container {
	return &resources.Container{
		Container: tau.Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
		},
	}
}

func ContainerToProto(container resources.Container) *api.Container {
	return &api.Container{
		Name:    container.Name,
		Image:   container.Image,
		Command: container.Command,
	}
}
