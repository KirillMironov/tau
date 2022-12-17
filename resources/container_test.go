package resources

import (
	"reflect"
	"testing"

	"github.com/KirillMironov/tau"
)

func TestContainerGob(t *testing.T) {
	t.Parallel()

	var (
		want = Container{
			Name:    "name",
			Image:   "image",
			Command: "command",
			status:  tau.Status{State: tau.StateSucceeded},
		}
		got Container
	)

	data, err := want.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal container: %v", err)
	}

	if err = got.UnmarshalBinary(data); err != nil {
		t.Fatalf("failed to unmarshal container: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
