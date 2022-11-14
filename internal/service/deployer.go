package service

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type Deployer struct {
	createCh <-chan tau.Resource
	removeCh <-chan tau.Resource
	runtime  tau.ContainerRuntime
	validate *validator.Validate
	logger   logger.Logger
}

func NewDeployer(create, remove <-chan tau.Resource, runtime tau.ContainerRuntime, logger logger.Logger) *Deployer {
	return &Deployer{
		createCh: create,
		removeCh: remove,
		runtime:  runtime,
		validate: validator.New(),
		logger:   logger,
	}
}

func (d Deployer) Start() {
	for {
		select {
		case resource := <-d.createCh:
			err := d.create(resource)
			if err != nil {
				d.logger.Error(err)
			}
		case resource := <-d.removeCh:
			err := d.remove(resource)
			if err != nil {
				d.logger.Error(err)
			}
		}
	}
}

func (d Deployer) create(resource tau.Resource) error {
	err := d.validate.Struct(resource)
	if err != nil {
		return err
	}

	err = resource.Create(d.runtime)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return nil
}

func (d Deployer) remove(resource tau.Resource) error {
	err := resource.Delete(d.runtime)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return nil
}
