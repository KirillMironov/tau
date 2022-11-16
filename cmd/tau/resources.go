package main

import (
	"os"

	"github.com/KirillMironov/tau/api"
	"github.com/urfave/cli/v2"
)

type resources struct {
	client api.ResourcesClient
}

func (r resources) create() *cli.Command {
	var tomlFile string

	return &cli.Command{
		Name:                   "create",
		Usage:                  "create a resource from a toml file",
		UsageText:              "tau create -f <file>",
		Category:               "resources",
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			data, err := os.ReadFile(tomlFile)
			if err != nil {
				return err
			}

			_, err = r.client.Create(ctx.Context, &api.Request{Data: data})
			return err
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "f",
				Usage:       "path to the toml file",
				Required:    true,
				Destination: &tomlFile,
			},
		},
	}
}

func (r resources) remove() *cli.Command {
	var tomlFile string

	return &cli.Command{
		Name:                   "remove",
		Usage:                  "remove a resource from a toml file",
		UsageText:              "tau remove -f <file>",
		Category:               "resources",
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			data, err := os.ReadFile(tomlFile)
			if err != nil {
				return err
			}

			_, err = r.client.Remove(ctx.Context, &api.Request{Data: data})
			return err
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "f",
				Usage:       "path to the toml file",
				Required:    true,
				Destination: &tomlFile,
			},
		},
	}
}
