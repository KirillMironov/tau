package resources

import "github.com/KirillMironov/tau/api"

const (
	fileFlag = "file"

	nameFlag    = "name"
	imageFlag   = "image"
	commandFlag = "command"
)

type Group struct {
	client api.ResourcesClient
}

func NewGroup(client api.ResourcesClient) *Group {
	return &Group{client: client}
}
