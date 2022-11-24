package service

import (
	"fmt"

	"github.com/KirillMironov/tau/resources"
	"github.com/KirillMironov/tau/runtimes"
)

type Deployer struct {
	runtime runtimes.ContainerRuntime
}

func NewDeployer(runtime runtimes.ContainerRuntime) *Deployer {
	return &Deployer{runtime: runtime}
}

func (d Deployer) Create(resource resources.Resource) error {
	err := resource.Validate()
	if err != nil {
		return err
	}

	err = resource.Create(d.runtime)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return nil
}

func (d Deployer) Remove(resource resources.Resource) error {
	err := resource.Validate()
	if err != nil {
		return err
	}

	err = resource.Remove(d.runtime)
	if err != nil {
		return fmt.Errorf("failed to remove resource: %w", err)
	}

	return nil
}
