//go:generate mockgen -destination=./../pkg/mock/tau.go -package=mock . ContainerRuntime

package runtimes

type Container struct {
	Name    string
	Image   string
	Command string
}

type ContainerRuntime interface {
	Start(Container) error
	Remove(containerName string) error
}
