//go:generate mockgen -destination=./pkg/mock/tau.go -package=mock . ContainerRuntime

package tau

import "time"

type ContainerState int

const (
	ContainerStateRunning ContainerState = iota + 1
	ContainerStateSucceeded
	ContainerStateFailed
)

type Container struct {
	Name    string
	Image   string
	Command string
}

type ContainerRuntime interface {
	// Start creates and starts a container.
	Start(Container) error
	// Stop stops a running container.
	// Runtime waits for the container to stop gracefully for the given timeout.
	// After the timeout, the container is killed.
	// If the container is not found, ErrContainerNotFound is returned.
	Stop(containerName string, timeout time.Duration) error
	// Remove removes a container.
	// Running containers are killed before removal.
	// If the container is not found, no error is returned.
	Remove(containerName string) error
	// State returns the state of a container.
	// If the container is not found, ErrContainerNotFound is returned.
	State(containerName string) (ContainerState, error)
}
