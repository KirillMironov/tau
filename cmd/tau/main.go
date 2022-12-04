package main

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/cli/resources"
	"github.com/KirillMironov/tau/pkg/cmdutil"
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
	var client = api.NewResourcesClient(conn)

	// CLI
	tau := &cobra.Command{
		Use:           "tau",
		Short:         "tau is a CLI tool for managing resources",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	tau.AddCommand(
		resources.NewCommand(client)...,
	)

	cmdutil.DisableFlagsSorting(tau)
	cmdutil.HideHelpCommand(tau)
	cmdutil.HideHelpFlags(tau)

	err = tau.Execute()
	if err != nil {
		cmdutil.Exit(1, err)
	}
}
