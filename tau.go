//go:generate mockgen -destination=mock/tau.go -package=mock . ContainerRuntime

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
		Validate() error
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
