package model

import "github.com/KirillMironov/tau/api"

type Container struct {
	Name    string `toml:"name"`
	Image   string `toml:"image"`
	Command string `toml:"command"`
}

func (c Container) Descriptor() *api.Descriptor {
	return &api.Descriptor{
		Name: c.Name,
		Kind: api.Kind_KIND_CONTAINER,
	}
}

func (c Container) ToProto() *api.Resource {
	return &api.Resource{
		Kind: &api.Resource_Container{
			Container: &api.Container{
				Name:    c.Name,
				Image:   c.Image,
				Command: c.Command,
			},
		},
	}
}
