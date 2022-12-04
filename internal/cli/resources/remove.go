package resources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/pkg/cmdutil"
	"github.com/KirillMironov/tau/pkg/tomlutil"
)

func remove(client api.ResourcesClient) *cobra.Command {
	var tomlFile string

	command := &cobra.Command{
		Use:     "remove -f <file>",
		Short:   "remove a resource",
		Long:    "remove a resource from a toml file",
		Example: "tau remove -f resource.toml",
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

			kind, err := protoconv.KindToProto(resource.Kind())
			if err != nil {
				return err
			}

			request := &api.RemoveRequest{
				Kind: kind,
				Name: resource.ID(),
			}

			_, err = client.Remove(cmd.Context(), request)
			return err
		},
	}

	command.Flags().StringVarP(&tomlFile, fileFlag, cmdutil.ShortFlag(fileFlag), "", "path to a toml file")

	_ = command.MarkFlagRequired(fileFlag)

	command.AddCommand(
		removeContainer(client),
		removePod(client),
	)

	return command
}

func removeContainer(client api.ResourcesClient) *cobra.Command {
	command := &cobra.Command{
		Use:     "container <name>",
		Short:   "remove a resource",
		Long:    "remove a container",
		Example: "tau remove container busybox",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]

			request := &api.RemoveRequest{
				Kind: api.Kind_KIND_CONTAINER,
				Name: containerName,
			}

			_, err := client.Remove(cmd.Context(), request)
			return err
		},
	}

	return command
}

func removePod(client api.ResourcesClient) *cobra.Command {
	command := &cobra.Command{
		Use:     "pod <name>",
		Short:   "remove a resource",
		Long:    "remove a pod",
		Example: "tau remove pod busybox",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			podName := args[0]

			request := &api.RemoveRequest{
				Kind: api.Kind_KIND_POD,
				Name: podName,
			}

			_, err := client.Remove(cmd.Context(), request)
			return err
		},
	}

	return command
}
