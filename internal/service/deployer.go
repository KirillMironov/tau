package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/KirillMironov/tau"
)

type Deployer struct {
	runtime  tau.ContainerRuntime
	validate *validator.Validate
}

func NewDeployer(runtime tau.ContainerRuntime) *Deployer {
	return &Deployer{
		runtime:  runtime,
		validate: validator.New(),
	}
}

func (d Deployer) Create(resource tau.Resource) error {
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

func (d Deployer) Remove(resource tau.Resource) error {
	err := resource.Delete(d.runtime)
	if err != nil {
		return fmt.Errorf("failed to remove resource: %w", err)
	}

	return nil
}
