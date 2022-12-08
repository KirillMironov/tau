package service

import (
	"fmt"

	"github.com/KirillMironov/tau/pkg/logger"
	"github.com/KirillMironov/tau/resources"
	"github.com/KirillMironov/tau/runtimes"
)

type Resources struct {
	createCh <-chan resources.Resource
	removeCh <-chan resources.Descriptor
	storage  storage
	runtime  runtimes.ContainerRuntime
	logger   logger.Logger
}

type storage interface {
	Put(resources.Resource) error
	Get(resources.Descriptor) (resources.Resource, error)
	Delete(resources.Descriptor) error
}

func NewResources(createCh <-chan resources.Resource, removeCh <-chan resources.Descriptor, storage storage, runtime runtimes.ContainerRuntime, logger logger.Logger) *Resources {
	return &Resources{
		createCh: createCh,
		removeCh: removeCh,
		storage:  storage,
		runtime:  runtime,
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
	err := r.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(r.runtime)
}

func (r Resources) remove(descriptor resources.Descriptor) error {
	resource, err := r.storage.Get(descriptor)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	err = r.storage.Delete(descriptor)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return resource.Remove(r.runtime)
}
