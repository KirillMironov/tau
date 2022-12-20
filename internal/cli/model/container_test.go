package model

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/KirillMironov/tau/api"
)

func TestContainer_ToProto(t *testing.T) {
	t.Parallel()

	var (
		container = &Container{
			Name:    "container",
			Image:   "image",
			Command: "command",
		}
		want = &api.Resource{
			Kind: &api.Resource_Container{
				Container: &api.Container{
					Name:    "container",
					Image:   "image",
					Command: "command",
				},
			},
		}
	)

	if got := container.ToProto(); !proto.Equal(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
