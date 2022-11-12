package tau

import (
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

type (
	Resource interface {
		Create(ContainerRuntime) error
		Delete(ContainerRuntime) error
	}

	ContainerRuntime interface {
		Start(Container) (containerId string, _ error)
		Remove(containerId string) error
	}
)

// Disable Podman bindings logging (https://github.com/containers/podman/issues/13504).
func init() {
	logrus.SetOutput(io.Discard)
}
