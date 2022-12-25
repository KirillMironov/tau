package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/internal/domain"
	"github.com/KirillMironov/tau/pkg/logger"
)

type ResourcesController struct {
	runtime      tau.ContainerRuntime
	storage      resourcesStorage
	logger       logger.Logger
	updateTicker *time.Ticker
	once         sync.Once
}

type resourcesStorage interface {
	Put(tau.Resource) error
	Get(tau.Descriptor) (tau.Resource, error)
	Delete(tau.Descriptor) error
	List() ([]tau.Resource, error)
	ListByKind(tau.Kind) ([]tau.Resource, error)
}

func NewResourcesController(runtime tau.ContainerRuntime, storage resourcesStorage, logger logger.Logger, updateInterval time.Duration) *ResourcesController {
	return &ResourcesController{
		runtime:      runtime,
		storage:      storage,
		logger:       logger,
		updateTicker: time.NewTicker(updateInterval),
	}
}

func (rc *ResourcesController) Start() {
	rc.once.Do(func() {
		defer rc.updateTicker.Stop()

		for range rc.updateTicker.C {
			if err := rc.updateStatuses(); err != nil {
				rc.logger.Errorf("failed to update statuses: %v", err)
			}
		}
	})
}

func (rc *ResourcesController) Create(resource tau.Resource) error {
	rc.logger.Debugf("creating resource %#v", resource)

	err := rc.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(rc.runtime)
}

func (rc *ResourcesController) Remove(descriptor tau.Descriptor) error {
	rc.logger.Debugf("removing resource %#v", descriptor)

	resource, err := rc.storage.Get(descriptor)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	err = rc.storage.Delete(descriptor)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return resource.Remove(rc.runtime)
}

func (rc *ResourcesController) Status(descriptor tau.Descriptor) (tau.State, []tau.StatusEntry, error) {
	rc.logger.Debugf("getting status of resource %#v", descriptor)

	resource, err := rc.storage.Get(descriptor)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource.State(), resource.Status(), nil
}

func (rc *ResourcesController) updateStatuses() error {
	resources, err := rc.storage.List()
	if err != nil && !errors.Is(err, domain.ErrNoResources) {
		return fmt.Errorf("failed to list resources: %w", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(resources))

	for _, resource := range resources {
		go func(resource tau.Resource) {
			defer wg.Done()

			if err := resource.UpdateStatus(rc.runtime); err != nil {
				if !errors.Is(err, tau.ErrContainerNotFound) {
					rc.logger.Errorf("failed to update status for resource %+v: %v", resource.Descriptor(), err)
				}
				return
			}

			if err := rc.storage.Put(resource); err != nil {
				rc.logger.Errorf("failed to put resource %+v: %v", resource.Descriptor(), err)
				return
			}
		}(resource)
	}

	wg.Wait()

	return nil
}
