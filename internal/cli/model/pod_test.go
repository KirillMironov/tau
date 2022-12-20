package model

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/KirillMironov/tau/api"
)

func TestPod_ToProto(t *testing.T) {
	t.Parallel()

	var (
		pod = &Pod{
			Name: "pod",
			Containers: []container{
				{Name: "container", Image: "image", Command: "command"},
				{Name: "container2", Image: "image2", Command: "command2"},
			},
		}
		want = &api.Resource{
			Kind: &api.Resource_Pod{
				Pod: &api.Pod{
					Name: "pod",
					Containers: []*api.Pod_Container{
						{Name: "container", Image: "image", Command: "command"},
						{Name: "container2", Image: "image2", Command: "command2"},
					},
				},
			},
		}
	)

	if got := pod.ToProto(); !proto.Equal(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
