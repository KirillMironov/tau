package model

import "github.com/KirillMironov/tau/api"

type Pod struct {
	Name       string      `toml:"name"`
	Containers []container `toml:"containers"`
}

type container struct {
	Name    string `toml:"name"`
	Image   string `toml:"image"`
	Command string `toml:"command"`
}

func (p Pod) Descriptor() *api.Descriptor {
	return &api.Descriptor{
		Name: p.Name,
		Kind: api.Kind_KIND_POD,
	}
}

func (p Pod) ToProto() *api.Resource {
	containers := make([]*api.Pod_Container, 0, len(p.Containers))

	for _, v := range p.Containers {
		containers = append(containers, &api.Pod_Container{
			Name:    v.Name,
			Image:   v.Image,
			Command: v.Command,
		})
	}

	return &api.Resource{
		Kind: &api.Resource_Pod{
			Pod: &api.Pod{
				Name:       p.Name,
				Containers: containers,
			},
		},
	}
}
