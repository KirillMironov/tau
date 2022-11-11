package tau

import (
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

// Disable Podman bindings logging (https://github.com/containers/podman/issues/13504).
func init() {
	logrus.SetOutput(io.Discard)
}
