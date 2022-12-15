package tomlutil

import (
	"reflect"
	"testing"

	"github.com/KirillMironov/tau/resources"
)

func TestUnmarshalByKind(t *testing.T) {
	t.Parallel()

	test := []struct {
		name    string
		blob    []byte
		want    resources.Resource
		wantErr bool
	}{
		{
			name: "container",
			blob: []byte(`
				kind = "container"
				name = "busybox-sleep"
				image = "docker.io/library/busybox:1.35.0"
				command = "sleep 500"
			`),
			want: &resources.Container{
				Name:    "busybox-sleep",
				Image:   "docker.io/library/busybox:1.35.0",
				Command: "sleep 500",
			},
			wantErr: false,
		},
		{
			name: "pod",
			blob: []byte(`
				kind = "pod"
				name = "busybox"
				
				[[containers]]
				name = "busybox-sleep"
				image = "docker.io/library/busybox:1.35.0"
				command = "sleep 1000"
				
				[[containers]]
				name = "busybox-sleep"
				image = "docker.io/library/busybox:1.35.0"
				command = "sleep 500"
			`),
			want: &resources.Pod{
				Name: "busybox",
				Containers: []resources.Container{
					{
						Name:    "busybox-sleep",
						Image:   "docker.io/library/busybox:1.35.0",
						Command: "sleep 1000",
					},
					{
						Name:    "busybox-sleep",
						Image:   "docker.io/library/busybox:1.35.0",
						Command: "sleep 500",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "unexpected kind",
			blob:    []byte(""),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range test {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := UnmarshalByKind(tc.blob)
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
