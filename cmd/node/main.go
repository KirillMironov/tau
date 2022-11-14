package main

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/internal/service"
	"github.com/KirillMironov/tau/runtime"
	"github.com/sirupsen/logrus"
)

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
	createCh := make(chan tau.Resource)
	removeCh := make(chan tau.Resource)

	deployer := service.NewDeployer(createCh, removeCh, podmanRuntime, logger)

	// Node
	go deployer.Start()
}
