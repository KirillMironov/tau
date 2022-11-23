//go:generate mockgen -destination=./../pkg/mock/tau.go -package=mock . Runtime

package runtimes

type Container struct {
	Name    string
	Image   string
	Command string
}

type Runtime interface {
	Start(Container) error
	Remove(containerName string) error
}
