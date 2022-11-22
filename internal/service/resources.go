package service

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/logger"
)

type Resources struct {
	createCh <-chan tau.Resource
	removeCh <-chan tau.Resource
	deployer deployer
	logger   logger.Logger
}

type deployer interface {
	Create(tau.Resource) error
	Remove(tau.Resource) error
}

func NewResources(createCh, removeCh <-chan tau.Resource, deployer deployer, logger logger.Logger) *Resources {
	return &Resources{
		createCh: createCh,
		removeCh: removeCh,
		deployer: deployer,
		logger:   logger,
	}
}

func (r Resources) Start() {
	for {
		select {
		case resource := <-r.createCh:
			r.logger.Debugf("creating resource %#v", resource)

			err := r.deployer.Create(resource)
			if err != nil {
				r.logger.Error(err)
			}
		case resource := <-r.removeCh:
			r.logger.Debugf("removing resource %#v", resource)

			err := r.deployer.Remove(resource)
			if err != nil {
				r.logger.Error(err)
			}
		}
	}
}
