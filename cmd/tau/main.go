package main

import (
	"os"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/pkg/cmdutil"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddress = ":5668"

func main() {
	// gRPC connection
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cmdutil.Exitf(1, "failed to connect to gRPC server %q: %v", serverAddress, err)
	}
	defer conn.Close()

	// DI
	var (
		resources = resources{client: api.NewResourcesClient(conn)}
	)

	// CLI
	app := &cli.App{
		Name:  "tau",
		Usage: "tau is a CLI tool for managing resources",
		Commands: []*cli.Command{
			resources.create(),
			resources.remove(),
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		cmdutil.Exit(1, err)
	}
}
