package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/cli/resources"
	"github.com/KirillMironov/tau/pkg/cmdutil"
	"github.com/KirillMironov/tau/pkg/cobrax"
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
		client = api.NewResourcesClient(conn)
		group  = resources.NewGroup(client)
	)

	// CLI
	tau := &cobrax.Command{
		Usage:       "tau",
		Description: "tau is a CLI tool for managing resources",
		Options: cobrax.CommandOptions{
			SilenceErrors:       true,
			SilenceUsage:        true,
			DisableDefaultCmd:   true,
			DisableFlagsSorting: true,
			HideHelpCommand:     true,
			HideHelpFlags:       true,
		},
		Subcommands: []*cobrax.Command{
			group.Create(),
			group.Get(),
			group.Remove(),
		},
	}

	err = tau.Execute()
	if err != nil {
		cmdutil.Exit(1, err)
	}
}
