package resources

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"

	"github.com/KirillMironov/tau/pkg/mock"
)

func TestPod_Create(t *testing.T) {
	t.Parallel()

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

	if err := pod.Create(runtime); err == nil {
		t.Fatalf("err = %v, want not nil", err)
	}
}

func TestPod_Remove(t *testing.T) {
	t.Parallel()

	pod, runtime := setup(t)

	runtime.
		EXPECT().
		Remove(gomock.Any()).
		Return(errors.New("error")).
		Times(2)

	err := pod.Remove(runtime)
	if err == nil {
		t.Fatalf("err = %v, want not nil", err)
	}

	e, ok := err.(*multierror.Error)
	if !ok {
		t.Fatalf("err type = %T, want *multierror.Error", err)
	}

	if len(e.Errors) != 2 {
		t.Fatalf("len(e.Errors) = %d, want 2", len(e.Errors))
	}
}

func TestPodGob(t *testing.T) {
	t.Parallel()

	var (
		want = Pod{
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
		got Pod
	)

	data, err := want.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal pod: %v", err)
	}

	if err = got.UnmarshalBinary(data); err != nil {
		t.Fatalf("failed to unmarshal pod: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %+v, want %+v", got, want)
	}
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
