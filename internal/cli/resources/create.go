package resources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/cli/model"
	"github.com/KirillMironov/tau/pkg/cobrax"
)

func (g Group) Create() *cobrax.Command {
	return &cobrax.Command{
		Usage:       "create -f <file>",
		Description: "create a resource from a toml file",
		Example:     "tau create -f resource.toml",
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

			resource, err := model.UnmarshalByKind(data)
			if err != nil {
				return err
			}

			_, err = g.client.Create(cmd.Context(), resource.ToProto())
			return err
		},
		Subcommands: []*cobrax.Command{
			{
				Usage:       "container --name <name> --image <image> --command [command]",
				Description: "create a container",
				Example:     "tau create container --name busybox --image docker.io/library/busybox:1.35.0",
				Args:        cobra.NoArgs,
				Flags: []cobrax.Flag{
					&cobrax.StringFlag{
						Name:     nameFlag,
						Usage:    "container name",
						Required: true,
					},
					&cobrax.StringFlag{
						Name:     imageFlag,
						Usage:    "container image",
						Required: true,
					},
					&cobrax.StringFlag{
						Name:  commandFlag,
						Usage: "command to run in the container",
					},
				},
				Action: func(cmd *cobra.Command, _ []string) error {
					name := cmd.Flag(nameFlag).Value.String()
					image := cmd.Flag(imageFlag).Value.String()
					command := cmd.Flag(commandFlag).Value.String()

					container := &api.Resource{
						Kind: &api.Resource_Container{
							Container: &api.Container{
								Name:    name,
								Image:   image,
								Command: command,
							},
						},
					}

					_, err := g.client.Create(cmd.Context(), container)
					return err
				},
			},
		},
	}
}
