package main

import (
	"net"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/internal/service"
	"github.com/KirillMironov/tau/internal/transport"
	"github.com/KirillMironov/tau/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const address = ":5668"

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Runtime
	podmanRuntime, err := runtime.NewPodman(runtime.PodmanRootlessSocket())
	if err != nil {
		logger.Fatal(err)
	}

	// DI
	var (
		createCh = make(chan tau.Resource)
		removeCh = make(chan tau.Resource)

		deployer  = service.NewDeployer(createCh, removeCh, podmanRuntime, logger)
		resources = transport.NewResources(createCh, removeCh)
	)

	// Node
	go deployer.Start()

	// gRPC server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal(err)
	}

	server := grpc.NewServer()
	server.RegisterService(&api.Resources_ServiceDesc, resources)

	err = server.Serve(listener)
	if err != nil {
		logger.Fatal(err)
	}
}
