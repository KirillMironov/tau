package resources

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/pkg/mock"
)

func TestPod_Create(t *testing.T) {
	var (
		pod = Pod{
			Containers: make([]tau.Container, 2),
		}

		ctrl    = gomock.NewController(t)
		runtime = mock.NewMockContainerRuntime(ctrl)
	)

	gomock.InOrder(
		runtime.EXPECT().Start(&tau.Container{}).Return(errors.New("error")).Times(1),
		runtime.EXPECT().Remove(gomock.Any()).Times(2),
	)

	err := pod.Create(runtime)
	require.Error(t, err)
}

func TestPod_Delete(t *testing.T) {
	var (
		pod = Pod{
			Containers: make([]tau.Container, 2),
		}

		ctrl    = gomock.NewController(t)
		runtime = mock.NewMockContainerRuntime(ctrl)
	)

	runtime.EXPECT().Remove(gomock.Any()).Return(errors.New("error")).Times(2)

	err := pod.Delete(runtime)
	require.Error(t, err)

	e, ok := err.(*multierror.Error)
	require.True(t, ok)
	require.Len(t, e.Errors, 2)
}
