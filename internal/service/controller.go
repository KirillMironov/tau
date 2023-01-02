package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/internal/domain"
)

type ResourcesController struct {
	runtime      tau.ContainerRuntime
	storage      resourcesStorage
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

func NewResourcesController(runtime tau.ContainerRuntime, storage resourcesStorage, updateInterval time.Duration) *ResourcesController {
	return &ResourcesController{
		runtime:      runtime,
		storage:      storage,
		updateTicker: time.NewTicker(updateInterval),
	}
}

func (rc *ResourcesController) Start() {
	rc.once.Do(func() {
		defer rc.updateTicker.Stop()

		for range rc.updateTicker.C {
			if err := rc.updateStatuses(); err != nil {
				slog.Error("failed to update resource statuses", err)
			}
		}
	})
}

func (rc *ResourcesController) Create(resource tau.Resource) error {
	slog.Debug("creating resource", slog.Any("descriptor", resource.Descriptor()))

	err := rc.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(rc.runtime)
}

func (rc *ResourcesController) Remove(descriptor tau.Descriptor) error {
	slog.Debug("removing resource", slog.Any("descriptor", descriptor))

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
	slog.Debug("getting status of resource", slog.Any("descriptor", descriptor))

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
					slog.Error("failed to update resource status", err, slog.Any("descriptor", resource.Descriptor()))
				}
				return
			}

			if err := rc.storage.Put(resource); err != nil {
				slog.Error("failed to put resource", err, slog.Any("descriptor", resource.Descriptor()))
				return
			}
		}(resource)
	}

	wg.Wait()

	return nil
}
