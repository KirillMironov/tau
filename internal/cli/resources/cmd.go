package resources

import (
	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
)

const (
	fileFlag = "file"

	nameFlag    = "name"
	imageFlag   = "image"
	commandFlag = "command"
)

func NewCommand(client api.ResourcesClient) []*cobra.Command {
	return []*cobra.Command{
		create(client).Build(),
		remove(client).Build(),
	}
}
