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

const updateStatusesInterval = time.Second * 2

type Controller struct {
	runtime tau.ContainerRuntime
	storage resourcesStorage
	logger  logger.Logger
}

func NewController(runtime tau.ContainerRuntime, storage resourcesStorage, logger logger.Logger) *Controller {
	return &Controller{
		runtime: runtime,
		storage: storage,
		logger:  logger,
	}
}

func (c Controller) Start() {
	ticker := time.NewTicker(updateStatusesInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := c.updateStatuses(); err != nil {
			c.logger.Errorf("failed to update statuses: %v", err)
		}
	}
}

func (c Controller) updateStatuses() error {
	resources, err := c.storage.List()
	if err != nil && !errors.Is(err, domain.ErrNoResources) {
		return fmt.Errorf("failed to list resources: %w", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(resources))

	for _, resource := range resources {
		go c.updateStatus(resource, wg)
	}

	wg.Wait()

	return nil
}

func (c Controller) updateStatus(resource tau.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := resource.UpdateStatus(c.runtime); err != nil {
		c.logger.Errorf("failed to update status for resource %+v: %v", resource.Descriptor(), err)
		return
	}

	if err := c.storage.Put(resource); err != nil {
		c.logger.Errorf("failed to put resource %+v: %v", resource.Descriptor(), err)
		return
	}
}
