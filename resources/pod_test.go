package resources

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	"github.com/KirillMironov/tau/pkg/mock"
)

func TestPod_Create(t *testing.T) {
	pod, runtime := setup(t)

	gomock.InOrder(
		runtime.
			EXPECT().
			Start(gomock.Any()).
			Return(errors.New("error")).
			Times(1),
		runtime.
			EXPECT().
			Remove(gomock.Any()).
			Times(2),
	)

	err := pod.Create(runtime)
	require.Error(t, err)
}

func TestPod_Remove(t *testing.T) {
	pod, runtime := setup(t)

	runtime.
		EXPECT().
		Remove(gomock.Any()).
		Return(errors.New("error")).
		Times(2)

	err := pod.Remove(runtime)
	require.Error(t, err)

	e, ok := err.(*multierror.Error)
	require.True(t, ok)
	require.Len(t, e.Errors, 2)
}

func TestPodGob(t *testing.T) {
	var (
		pod = Pod{
			Name: "name",
			Containers: []Container{
				{
					Name:    "name-1",
					Image:   "image-1",
					Command: "command-2",
					status:  Status{State: StateRunning},
				},
				{
					Name:   "name-2",
					Image:  "image-2",
					status: Status{State: StatePending},
				},
			},
			status: Status{State: StateSucceeded},
		}
		target Pod
	)

	data, err := pod.MarshalBinary()
	require.NoError(t, err)

	err = target.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, pod, target)
}

func setup(t *testing.T) (Pod, *mock.MockContainerRuntime) {
	t.Helper()

	var (
		pod = Pod{
			Name: "pod",
			Containers: []Container{
				{Name: "1", Image: "image"},
				{Name: "2", Image: "image"},
			},
		}
		ctrl    = gomock.NewController(t)
		runtime = mock.NewMockContainerRuntime(ctrl)
	)

	return pod, runtime
}
