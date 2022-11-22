package service

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/logger"
)

type Resources struct {
	createCh <-chan tau.Resource
	removeCh <-chan tau.Resource
	storage  storage
	deployer deployer
	logger   logger.Logger
}

type (
	storage interface {
		Create(tau.Resource) error
		GetById(id string) (tau.Resource, error)
		Delete(id string) error
	}
	deployer interface {
		Create(tau.Resource) error
		Remove(tau.Resource) error
	}
)

func NewResources(createCh, removeCh <-chan tau.Resource, storage storage, deployer deployer,
	logger logger.Logger) *Resources {
	return &Resources{
		createCh: createCh,
		removeCh: removeCh,
		storage:  storage,
		deployer: deployer,
		logger:   logger,
	}
}

func (r Resources) Start() {
	for {
		select {
		case resource := <-r.createCh:
			r.logger.Debugf("creating resource %#v", resource)

			err := r.create(resource)
			if err != nil {
				r.logger.Error(err)
			}
		case resource := <-r.removeCh:
			r.logger.Debugf("removing resource %#v", resource)

			err := r.remove(resource)
			if err != nil {
				r.logger.Error(err)
			}
		}
	}
}

func (r Resources) create(resource tau.Resource) error {
	err := r.storage.Create(resource)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return r.deployer.Create(resource)
}

func (r Resources) remove(resource tau.Resource) error {
	err := r.storage.Delete(resource.Id())
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return r.deployer.Remove(resource)
}
