package protoconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func TestResourceFromProto(t *testing.T) {
	resource, protoResource := testResources()

	converted, err := ResourceFromProto(protoResource)
	require.NoError(t, err)
	assert.Equal(t, resource, converted)

	converted, err = ResourceFromProto(nil)
	require.Error(t, err)
	assert.Nil(t, converted)
}

func TestResourceToProto(t *testing.T) {
	resource, protoResource := testResources()

	converted, err := ResourceToProto(resource)
	require.NoError(t, err)
	assert.Equal(t, protoResource, converted)

	converted, err = ResourceToProto(&resource)
	require.NoError(t, err)
	assert.Equal(t, protoResource, converted)

	converted, err = ResourceToProto(nil)
	require.Error(t, err)
	assert.Nil(t, converted)
}

func testResources() (resources.Pod, *api.Resource) {
	var (
		resource = resources.Pod{
			Name: "pod",
			Containers: []tau.Container{
				{Image: "image", Command: "command"},
				{Image: "image2", Command: "command2"},
			},
		}
		protoResource = &api.Resource{
			Kind: &api.Resource_Pod{
				Pod: &api.Pod{
					Name: "pod",
					Containers: []*api.Container{
						{Image: "image", Command: "command"},
						{Image: "image2", Command: "command2"},
					},
				},
			},
		}
	)

	return resource, protoResource
}
