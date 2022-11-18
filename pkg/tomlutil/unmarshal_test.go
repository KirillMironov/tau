package tomlutil

import (
	"testing"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalByKind(t *testing.T) {
	var (
		blob = []byte(`
			kind = "pod"
			name = "busybox"

			[[containers]]
			image = "docker.io/library/busybox:latest"
			command = "sleep 1000"

			[[containers]]
			image = "docker.io/library/busybox:latest"
			command = "sleep 500"
		`)
		expected = &resources.Pod{
			Name: "busybox",
			Containers: []tau.Container{
				{
					Image:   "docker.io/library/busybox:latest",
					Command: "sleep 1000",
				},
				{
					Image:   "docker.io/library/busybox:latest",
					Command: "sleep 500",
				},
			},
		}
	)

	resource, err := UnmarshalByKind(blob)
	require.NoError(t, err)
	assert.Equal(t, expected, resource)

	resource, err = UnmarshalByKind(nil)
	require.Error(t, err)
	assert.Empty(t, resource)
}
