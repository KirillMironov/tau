package service

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/logger"
)

type Resources struct {
	storage storage
	runtime tau.ContainerRuntime
	logger  logger.Logger
}

type storage interface {
	Put(tau.Resource) error
	Get(tau.Descriptor) (tau.Resource, error)
	Delete(tau.Descriptor) error
}

func NewResources(storage storage, runtime tau.ContainerRuntime, logger logger.Logger) *Resources {
	return &Resources{
		storage: storage,
		runtime: runtime,
		logger:  logger,
	}
}

func (r Resources) Create(resource tau.Resource) error {
	r.logger.Debugf("creating resource %#v", resource)

	resource.SetState(tau.StatePending)

	err := r.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(r.runtime)
}

func (r Resources) Get(descriptor tau.Descriptor) (tau.Status, error) {
	r.logger.Debugf("getting resource %#v", descriptor)

	resource, err := r.storage.Get(descriptor)
	if err != nil {
		return tau.Status{}, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource.Status(), nil
}

func (r Resources) Remove(descriptor tau.Descriptor) error {
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
