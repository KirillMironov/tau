package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainerGob(t *testing.T) {
	var (
		container = Container{
			Name:    "name",
			Image:   "image",
			Command: "command",
			status:  Status{State: StateSucceeded},
		}
		target Container
	)

	data, err := container.MarshalBinary()
	require.NoError(t, err)

	err = target.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, container, target)
}
