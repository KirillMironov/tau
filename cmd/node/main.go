package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/client"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/resources"
	"github.com/KirillMironov/tau/pkg/logger"
	"github.com/KirillMironov/tau/runtimes"
)

const address = ":5668"

func main() {
	// Logger
	slog.SetDefault(logger.New(slog.LevelDebug))

	// BoltDB
	db, err := bolt.Open("tau.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		fatal(err)
	}
	defer db.Close()

	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fatal(err)
	}
	defer dockerClient.Close()

	_, err = dockerClient.Ping(context.Background())
	if err != nil {
		fatal(err)
	}

	// DI
	resourcesStorage, err := resources.NewStorage(db)
	if err != nil {
		fatal(err)
	}

	var (
		runtime = runtimes.NewDocker(dockerClient)

		resourcesService = resources.NewService(runtime, resourcesStorage, time.Second)
		resourcesHandler = resources.NewHandler(resourcesService)
	)

	// gRPC server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fatal(err)
	}

	server := grpc.NewServer()
	server.RegisterService(&api.Resources_ServiceDesc, resourcesHandler)

	err = server.Serve(listener)
	if err != nil {
		fatal(err)
	}
}

func fatal(err error, args ...any) {
	slog.Error("", err, args)
	os.Exit(1)
}
