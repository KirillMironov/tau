package protoconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func TestResourceFromProto_Container(t *testing.T) {
	container, _, protoContainer, _ := testResources()

	resource, err := ResourceFromProto(protoContainer)
	require.NoError(t, err)
	assert.Equal(t, container, resource)

	resource, err = ResourceFromProto(nil)
	require.Error(t, err)
	assert.Nil(t, resource)
}

func TestResourceFromProto_Pod(t *testing.T) {
	_, pod, _, protoPod := testResources()

	resource, err := ResourceFromProto(protoPod)
	require.NoError(t, err)
	assert.Equal(t, pod, resource)

	resource, err = ResourceFromProto(nil)
	require.Error(t, err)
	assert.Nil(t, resource)
}

func TestResourceToProto_Container(t *testing.T) {
	container, _, protoContainer, _ := testResources()

	resource, err := ResourceToProto(container)
	require.NoError(t, err)
	assert.Equal(t, protoContainer, resource)

	resource, err = ResourceToProto(container)
	require.NoError(t, err)
	assert.Equal(t, protoContainer, resource)

	resource, err = ResourceToProto(nil)
	require.Error(t, err)
	assert.Nil(t, resource)
}

func TestResourceToProto_Pod(t *testing.T) {
	_, pod, _, protoPod := testResources()

	resource, err := ResourceToProto(pod)
	require.NoError(t, err)
	assert.Equal(t, protoPod, resource)

	resource, err = ResourceToProto(pod)
	require.NoError(t, err)
	assert.Equal(t, protoPod, resource)

	resource, err = ResourceToProto(nil)
	require.Error(t, err)
	assert.Nil(t, resource)
}

func testResources() (*resources.Container, *resources.Pod, *api.Resource, *api.Resource) {
	var (
		container = &resources.Container{
			Name:    "container",
			Image:   "image",
			Command: "command",
		}
		pod = &resources.Pod{
			Name: "pod",
			Containers: []resources.Container{
				{Name: "container", Image: "image", Command: "command"},
				{Name: "container2", Image: "image2", Command: "command2"},
			},
		}
		protoContainer = &api.Resource{
			Kind: &api.Resource_Container{
				Container: &api.Container{
					Name:    "container",
					Image:   "image",
					Command: "command",
				},
			},
		}
		protoPod = &api.Resource{
			Kind: &api.Resource_Pod{
				Pod: &api.Pod{
					Name: "pod",
					Containers: []*api.Container{
						{Name: "container", Image: "image", Command: "command"},
						{Name: "container2", Image: "image2", Command: "command2"},
					},
				},
			},
		}
	)

	return container, pod, protoContainer, protoPod
}
