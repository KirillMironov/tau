package main

import (
	"context"
	"net"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/service"
	"github.com/KirillMironov/tau/internal/storage"
	"github.com/KirillMironov/tau/internal/transport"
	"github.com/KirillMironov/tau/runtimes"
)

const address = ":5668"

func main() {
	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// BoltDB
	db, err := bolt.Open("tau.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal(err)
	}
	defer dockerClient.Close()

	_, err = dockerClient.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	// DI
	resourcesStorage, err := storage.NewResources(db)
	if err != nil {
		logger.Fatal(err)
	}

	var (
		runtime = runtimes.NewDocker(dockerClient)

		resourcesController = service.NewResourcesController(runtime, resourcesStorage, logger, time.Second)
		resourcesHandler    = transport.NewResources(resourcesController)
	)

	// Services
	go resourcesController.Start()

	// gRPC server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal(err)
	}

	server := grpc.NewServer()
	server.RegisterService(&api.Resources_ServiceDesc, resourcesHandler)

	err = server.Serve(listener)
	if err != nil {
		logger.Fatal(err)
	}
}
