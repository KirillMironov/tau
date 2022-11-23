package tomlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KirillMironov/tau/resources"
)

func TestUnmarshalByKind_Container(t *testing.T) {
	var (
		blob = []byte(`
			kind = "container"
			name = "busybox-sleep"
			image = "docker.io/library/busybox:1.35.0"
			command = "sleep 500"
		`)
		container = &resources.Container{
			Name:    "busybox-sleep",
			Image:   "docker.io/library/busybox:1.35.0",
			Command: "sleep 500",
		}
	)

	resource, err := UnmarshalByKind(blob)
	require.NoError(t, err)
	assert.Equal(t, container, resource)
}

func TestUnmarshalByKind_Pod(t *testing.T) {
	var (
		blob = []byte(`
			kind = "pod"
			name = "busybox"

			[[containers]]
			name = "busybox-sleep"
			image = "docker.io/library/busybox:1.35.0"
			command = "sleep 1000"

			[[containers]]
			name = "busybox-sleep"
			image = "docker.io/library/busybox:1.35.0"
			command = "sleep 500"
		`)
		pod = &resources.Pod{
			Name: "busybox",
			Containers: []resources.Container{
				{
					Name:    "busybox-sleep",
					Image:   "docker.io/library/busybox:1.35.0",
					Command: "sleep 1000",
				},
				{
					Name:    "busybox-sleep",
					Image:   "docker.io/library/busybox:1.35.0",
					Command: "sleep 500",
				},
			},
		}
	)

	resource, err := UnmarshalByKind(blob)
	require.NoError(t, err)
	assert.Equal(t, pod, resource)
}

func TestUnmarshalByKind(t *testing.T) {
	resource, err := UnmarshalByKind(nil)
	require.Error(t, err)
	assert.Empty(t, resource)
}
