package resources

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/KirillMironov/tau"
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

	e1 := errors.New("error 1")
	e2 := errors.New("error 2")

	pod, runtime := setup(t)

	runtime.
		EXPECT().
		Remove(gomock.Any()).
		Return(e1).
		Times(1)

	runtime.
		EXPECT().
		Remove(gomock.Any()).
		Return(e2).
		Times(1)

	err := pod.Remove(runtime)
	if err == nil {
		t.Fatalf("err = %v, want not nil", err)
	}

	if !errors.Is(err, e1) {
		t.Fatalf("err = %v, want %v", err, e1)
	}

	if !errors.Is(err, e2) {
		t.Fatalf("err = %v, want %v", err, e2)
	}
}

func TestPodGob(t *testing.T) {
	t.Parallel()

	var (
		want = Pod{
			Name: "name",
			Containers: []tau.Container{
				{
					Name:    "name-1",
					Image:   "image-1",
					Command: "command-2",
				},
				{
					Name:  "name-2",
					Image: "image-2",
				},
			},
			state: tau.StateRunning,
			status: podStatus{
				CreatedAt: time.Now().Unix(),
			},
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
			Containers: []tau.Container{
				{Name: "1", Image: "image"},
				{Name: "2", Image: "image"},
			},
		}
		ctrl    = gomock.NewController(t)
		runtime = mock.NewMockContainerRuntime(ctrl)
	)

	return pod, runtime
}
