package protoconv

import (
	"reflect"
	"testing"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

var (
	container = &resources.Container{
		Container: tau.Container{
			Name:    "container",
			Image:   "image",
			Command: "command",
		},
	}
	pod = &resources.Pod{
		Name: "pod",
		Containers: []tau.Container{
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

func TestResourceFromProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		proto   *api.Resource
		want    tau.Resource
		wantErr bool
	}{
		{
			name:    "container",
			proto:   protoContainer,
			want:    container,
			wantErr: false,
		},
		{
			name:    "pod",
			proto:   protoPod,
			want:    pod,
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

func TestResourceToProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resource tau.Resource
		want     *api.Resource
		wantErr  bool
	}{
		{
			name:     "container",
			resource: container,
			want:     protoContainer,
			wantErr:  false,
		},
		{
			name:     "pod",
			resource: pod,
			want:     protoPod,
			wantErr:  false,
		},
		{
			name:     "unexpected resource type",
			resource: nil,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ResourceToProto(tc.resource)
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
