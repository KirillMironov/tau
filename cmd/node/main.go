package main

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/service"
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

	// Runtime
	runtime, err := runtimes.NewPodman(runtimes.PodmanRootlessSocket())
	if err != nil {
		logger.Fatal(err)
	}

	// DI
	var (
		createCh = make(chan tau.Resource)
		removeCh = make(chan tau.Resource)

		deployer  = service.NewDeployer(runtime)
		resources = service.NewResources(createCh, removeCh, deployer, logger)

		resourcesServer = transport.NewResources(createCh, removeCh)
	)

	// Node
	go resources.Start()

	// gRPC server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal(err)
	}

	server := grpc.NewServer()
	server.RegisterService(&api.Resources_ServiceDesc, resourcesServer)

	err = server.Serve(listener)
	if err != nil {
		logger.Fatal(err)
	}
}
