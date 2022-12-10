package service

import (
	"fmt"

	"github.com/KirillMironov/tau/pkg/logger"
	"github.com/KirillMironov/tau/resources"
	"github.com/KirillMironov/tau/runtimes"
)

type Resources struct {
	storage storage
	runtime runtimes.ContainerRuntime
	logger  logger.Logger
}

type storage interface {
	Put(resources.Resource) error
	Get(resources.Descriptor) (resources.Resource, error)
	Delete(resources.Descriptor) error
}

func NewResources(storage storage, runtime runtimes.ContainerRuntime, logger logger.Logger) *Resources {
	return &Resources{
		storage: storage,
		runtime: runtime,
		logger:  logger,
	}
}

func (r Resources) Create(resource resources.Resource) error {
	r.logger.Debugf("creating resource %#v", resource)

	resource.SetState(resources.StatePending)

	err := r.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(r.runtime)
}

func (r Resources) Get(descriptor resources.Descriptor) (resources.Status, error) {
	r.logger.Debugf("getting resource %#v", descriptor)

	resource, err := r.storage.Get(descriptor)
	if err != nil {
		return resources.Status{}, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource.Status(), nil
}

func (r Resources) Remove(descriptor resources.Descriptor) error {
	r.logger.Debugf("removing resource %#v", descriptor)

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
