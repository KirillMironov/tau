package tau

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Disable Podman bindings logging (https://github.com/containers/podman/issues/13504).
func init() {
	logrus.SetOutput(io.Discard)
}
