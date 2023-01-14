package resources

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/KirillMironov/tau"
)

var errNoResources = errors.New("no resources")

type Service struct {
	runtime        tau.ContainerRuntime
	storage        storage
	updateInterval time.Duration
}

type storage interface {
	Put(tau.Resource) error
	Get(tau.Descriptor) (tau.Resource, error)
	Delete(tau.Descriptor) error
	List() ([]tau.Resource, error)
	ListByKind(tau.Kind) ([]tau.Resource, error)
}

func NewService(runtime tau.ContainerRuntime, storage storage, updateInterval time.Duration) *Service {
	s := &Service{
		runtime:        runtime,
		storage:        storage,
		updateInterval: updateInterval,
	}

	go s.updateStatuses()

	return s
}

func (s *Service) Create(resource tau.Resource) error {
	slog.Debug("creating resource", slog.Any("descriptor", resource.Descriptor()))

	err := s.storage.Put(resource)
	if err != nil {
		return fmt.Errorf("failed to put resource: %w", err)
	}

	return resource.Create(s.runtime)
}

func (s *Service) Remove(descriptor tau.Descriptor) error {
	slog.Debug("removing resource", slog.Any("descriptor", descriptor))

	resource, err := s.storage.Get(descriptor)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	err = s.storage.Delete(descriptor)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return resource.Remove(s.runtime)
}

func (s *Service) Status(descriptor tau.Descriptor) (tau.State, []tau.StatusEntry, error) {
	slog.Debug("getting status of resource", slog.Any("descriptor", descriptor))

	resource, err := s.storage.Get(descriptor)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource.State(), resource.Status(), nil
}

func (s *Service) updateStatuses() {
	ticker := time.NewTicker(s.updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		resources, err := s.storage.List()
		if err != nil && !errors.Is(err, errNoResources) {
			slog.Error("failed to list resources", err)
			continue
		}

		wg := new(sync.WaitGroup)
		wg.Add(len(resources))

		for _, resource := range resources {
			go func(resource tau.Resource) {
				defer wg.Done()

				if err := resource.UpdateStatus(s.runtime); err != nil {
					if !errors.Is(err, tau.ErrContainerNotFound) {
						slog.Error("failed to update resource status", err, slog.Any("descriptor", resource.Descriptor()))
					}
					return
				}

				if err := s.storage.Put(resource); err != nil {
					slog.Error("failed to put resource", err, slog.Any("descriptor", resource.Descriptor()))
					return
				}
			}(resource)
		}

		wg.Wait()
	}
}
