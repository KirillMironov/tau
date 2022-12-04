package resources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/pkg/cmdutil"
	"github.com/KirillMironov/tau/pkg/tomlutil"
)

func create(client api.ResourcesClient) *cobra.Command {
	var tomlFile string

	command := &cobra.Command{
		Use:     "create -f <file>",
		Short:   "create a resource",
		Long:    "create a resource from a toml file",
		Example: "tau create -f resource.toml",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := os.ReadFile(tomlFile)
			if err != nil {
				return err
			}

			resource, err := tomlutil.UnmarshalByKind(data)
			if err != nil {
				return err
			}

			protoResource, err := protoconv.ResourceToProto(resource)
			if err != nil {
				return err
			}

			_, err = client.Create(cmd.Context(), protoResource)
			return err
		},
	}

	command.Flags().StringVarP(&tomlFile, fileFlag, cmdutil.ShortFlag(fileFlag), "", "path to a toml file")

	_ = command.MarkFlagRequired(fileFlag)

	command.AddCommand(createContainer(client))

	return command
}

func createContainer(client api.ResourcesClient) *cobra.Command {
	var (
		name             string
		image            string
		containerCommand string
	)

	command := &cobra.Command{
		Use:     "container --name <name> --image <image> --command [command]",
		Short:   "create a resource",
		Long:    "create a container",
		Example: "tau create container --name busybox --image docker.io/library/busybox:1.35.0",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			container := &api.Resource{
				Kind: &api.Resource_Container{
					Container: &api.Container{
						Name:    name,
						Image:   image,
						Command: containerCommand,
					},
				},
			}

			_, err := client.Create(cmd.Context(), container)
			return err
		},
	}

	command.Flags().StringVarP(&name, nameFlag, cmdutil.ShortFlag(nameFlag), "", "container name")
	command.Flags().StringVarP(&image, imageFlag, cmdutil.ShortFlag(imageFlag), "", "container image")
	command.Flags().StringVarP(&containerCommand, commandFlag, cmdutil.ShortFlag(commandFlag), "", "command to run in the container")

	_ = command.MarkFlagRequired(nameFlag)
	_ = command.MarkFlagRequired(imageFlag)

	return command
}
