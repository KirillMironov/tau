package service

import (
	"fmt"

	"github.com/KirillMironov/tau/pkg/logger"
	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	createCh <-chan resources.Resource
	removeCh <-chan resources.Resource
	storage  storage
	deployer deployer
	logger   logger.Logger
}

type (
	storage interface {
		Create(resources.Resource) error
		GetById(id string) (resources.Resource, error)
		Delete(id string) error
	}
	deployer interface {
		Create(resources.Resource) error
		Remove(resources.Resource) error
	}
)

func NewResources(createCh, removeCh <-chan resources.Resource, storage storage, deployer deployer,
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

func (r Resources) create(resource resources.Resource) error {
	err := r.storage.Create(resource)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return r.deployer.Create(resource)
}

func (r Resources) remove(resource resources.Resource) error {
	err := r.storage.Delete(resource.ID())
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return r.deployer.Remove(resource)
}
