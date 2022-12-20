package service

import (
	"fmt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/logger"
)

type Resources struct {
	runtime tau.ContainerRuntime
	storage resourcesStorage
	logger  logger.Logger
}

type resourcesStorage interface {
	Put(tau.Resource) error
	Get(tau.Descriptor) (tau.Resource, error)
	Delete(tau.Descriptor) error
	List() ([]tau.Resource, error)
	ListByKind(tau.Kind) ([]tau.Resource, error)
}

func NewResources(runtime tau.ContainerRuntime, storage resourcesStorage, logger logger.Logger) *Resources {
	return &Resources{
		runtime: runtime,
		storage: storage,
		logger:  logger,
	}
}

func (r Resources) Create(resource tau.Resource) error {
	r.logger.Debugf("creating resource %#v", resource)

	err := r.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(r.runtime)
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

func (r Resources) Status(descriptor tau.Descriptor) (tau.State, []tau.StatusEntry, error) {
	r.logger.Debugf("getting status of resource %#v", descriptor)

	resource, err := r.storage.Get(descriptor)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource.State(), resource.Status(), nil
}
