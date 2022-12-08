package resources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/pkg/cobrax"
	"github.com/KirillMironov/tau/pkg/tomlutil"
)

func (g Group) Remove() *cobrax.Command {
	return &cobrax.Command{
		Usage:       "remove -f <file>",
		Description: "remove a resource from a toml file",
		Example:     "tau remove -f resource.toml",
		Args:        cobra.NoArgs,
		Flags: []cobrax.Flag{
			&cobrax.StringFlag{
				Name:     fileFlag,
				Usage:    "path to a toml file",
				Required: true,
			},
		},
		Action: func(cmd *cobra.Command, _ []string) error {
			tomlFile := cmd.Flag(fileFlag).Value.String()

			data, err := os.ReadFile(tomlFile)
			if err != nil {
				return err
			}

			resource, err := tomlutil.UnmarshalByKind(data)
			if err != nil {
				return err
			}

			descriptor := resource.Descriptor()

			kind, err := protoconv.KindToProto(descriptor.Kind)
			if err != nil {
				return err
			}

			request := &api.Descriptor{
				Name: descriptor.Name,
				Kind: kind,
			}

			_, err = g.client.Remove(cmd.Context(), request)
			return err
		},
		Subcommands: []*cobrax.Command{
			{
				Usage:       "container <name>",
				Description: "remove a container",
				Example:     "tau remove container busybox",
				Args:        cobra.ExactArgs(1),
				Action: func(cmd *cobra.Command, args []string) error {
					containerName := args[0]

					request := &api.Descriptor{
						Name: containerName,
						Kind: api.Kind_KIND_CONTAINER,
					}

					_, err := g.client.Remove(cmd.Context(), request)
					return err
				},
			},
			{
				Usage:       "pod <name>",
				Description: "remove a pod",
				Example:     "tau remove pod busybox",
				Args:        cobra.ExactArgs(1),
				Action: func(cmd *cobra.Command, args []string) error {
					podName := args[0]

					request := &api.Descriptor{
						Name: podName,
						Kind: api.Kind_KIND_POD,
					}

					_, err := g.client.Remove(cmd.Context(), request)
					return err
				},
			},
		},
	}
}
