package protoconv

import (
	"reflect"
	"testing"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func TestResourceFromProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		proto   *api.Resource
		want    tau.Resource
		wantErr bool
	}{
		{
			name: "container",
			proto: &api.Resource{
				Kind: &api.Resource_Container{
					Container: &api.Container{
						Name:    "container",
						Image:   "image",
						Command: "command",
					},
				},
			},
			want: &resources.Container{
				Container: tau.Container{
					Name:    "container",
					Image:   "image",
					Command: "command",
				},
			},
			wantErr: false,
		},
		{
			name: "pod",
			proto: &api.Resource{
				Kind: &api.Resource_Pod{
					Pod: &api.Pod{
						Name: "pod",
						Containers: []*api.Pod_Container{
							{Name: "container", Image: "image", Command: "command"},
							{Name: "container2", Image: "image2", Command: "command2"},
						},
					},
				},
			},
			want: &resources.Pod{
				Name: "pod",
				Containers: []tau.Container{
					{Name: "container", Image: "image", Command: "command"},
					{Name: "container2", Image: "image2", Command: "command2"},
				},
			},
			wantErr: false,
		},
		{
			name:    "unexpected resource type",
			proto:   &api.Resource{Kind: nil},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ResourceFromProto(tc.proto)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Errorf("err = %v, want error presence = %v", err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got = %+v, want = %+v", got, tc.want)
			}
		})
	}
}
