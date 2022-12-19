package main

import (
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
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Fatal(err)
	}
	defer dockerClient.Close()

	// DI
	resourcesStorage, err := storage.NewResources(db)
	if err != nil {
		logger.Fatal(err)
	}

	var (
		runtime = runtimes.NewDocker(dockerClient)

		controller       = service.NewController(runtime, resourcesStorage, logger)
		resourcesService = service.NewResources(runtime, resourcesStorage, logger)
		resourcesHandler = transport.NewResources(resourcesService)
	)

	// Services
	go controller.Start()

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
