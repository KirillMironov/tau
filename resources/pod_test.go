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

func setup(t *testing.T) (Pod, *mock.MockContainerRuntime) {
	t.Helper()

	var (
		pod     = Pod{Containers: make([]Container, 2)}
		ctrl    = gomock.NewController(t)
		runtime = mock.NewMockContainerRuntime(ctrl)
	)

	return pod, runtime
}
