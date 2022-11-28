package main

import (
	"net"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/service"
	"github.com/KirillMironov/tau/internal/storage"
	"github.com/KirillMironov/tau/internal/transport"
	"github.com/KirillMironov/tau/resources"
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

	// BadgerDB
	db, err := badger.Open(badger.DefaultOptions("./badger.db").WithLogger(nil))
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Runtime
	podman, err := runtimes.NewPodman(runtimes.PodmanRootlessSocket())
	if err != nil {
		logger.Fatal(err)
	}

	// DI
	var (
		createCh = make(chan resources.Resource)
		removeCh = make(chan resources.Resource)

		resourcesStorage = storage.NewResources(db)
		resourcesService = service.NewResources(createCh, removeCh, resourcesStorage, podman, logger)
		resourcesHandler = transport.NewResources(createCh, removeCh)
	)

	// Node
	go resourcesService.Start()

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
