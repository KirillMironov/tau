package tau

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPod(t *testing.T) {
	var (
		blob = []byte(`
			kind = "pod"
			name = "busybox"
			
			[[containers]]
			image = "docker.io/library/busybox:latest"
			command = ["sleep", "1000"]
			
			[[containers]]
			image = "docker.io/library/busybox:latest"
			command = ["sleep", "500"]
		`)
		expected = Pod{
			Kind: "pod",
			Name: "busybox",
			Containers: []Container{
				{
					Image:   "docker.io/library/busybox:latest",
					Command: []string{"sleep", "1000"},
				},
				{
					Image:   "docker.io/library/busybox:latest",
					Command: []string{"sleep", "500"},
				},
			},
		}
	)

	pod, err := NewPod(blob)
	require.NoError(t, err)
	assert.Equal(t, expected, pod)

	pod, err = NewPod(nil)
	require.Error(t, err)
	assert.Empty(t, pod)
}
