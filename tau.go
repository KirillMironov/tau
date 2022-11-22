//go:generate mockgen -destination=pkg/mock/tau.go -package=mock . ContainerRuntime

package tau

import (
	"io"

	"github.com/sirupsen/logrus"
)

type (
	Resource interface {
		Id() string
		Create(ContainerRuntime) error
		Delete(ContainerRuntime) error
	}

	ContainerRuntime interface {
		Start(*Container) error
		Remove(containerId string) error
	}
)

// Disable Podman bindings logging (https://github.com/containers/podman/issues/13504).
func init() {
	logrus.SetOutput(io.Discard)
}
