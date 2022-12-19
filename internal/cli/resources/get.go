package resources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/api/protoconv"
	"github.com/KirillMironov/tau/pkg/cobrax"
	"github.com/KirillMironov/tau/pkg/tomlutil"
)

func (g Group) Get() *cobrax.Command {
	return &cobrax.Command{
		Usage:       "get -f <file>",
		Description: "get a resource status",
		Example:     "tau get -f resource.toml",
		Args:        cobra.NoArgs,
		Flags: []cobrax.Flag{
			&cobrax.StringFlag{
				Name:     fileFlag,
				Usage:    "path to a toml file",
				Required: true,
			},
		},
		Action: func(cmd *cobra.Command, args []string) error {
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

			status, err := g.client.Status(cmd.Context(), request)
			if err != nil {
				return err
			}

			cmd.Println(status)

			return nil
		},
		Subcommands: []*cobrax.Command{
			{
				Usage:       "container <name>",
				Description: "get a container status",
				Example:     "tau get container busybox",
				Args:        cobra.ExactArgs(1),
				Action: func(cmd *cobra.Command, args []string) error {
					containerName := args[0]

					request := &api.Descriptor{
						Name: containerName,
						Kind: api.Kind_KIND_CONTAINER,
					}

					status, err := g.client.Status(cmd.Context(), request)
					if err != nil {
						return err
					}

					cmd.Println(status)

					return nil
				},
			},
			{
				Usage:       "pod <name>",
				Description: "get a pod status",
				Example:     "tau get pod busybox",
				Args:        cobra.ExactArgs(1),
				Action: func(cmd *cobra.Command, args []string) error {
					podName := args[0]

					request := &api.Descriptor{
						Name: podName,
						Kind: api.Kind_KIND_POD,
					}

					status, err := g.client.Status(cmd.Context(), request)
					if err != nil {
						return err
					}

					cmd.Println(status)

					return nil
				},
			},
		},
	}
}
