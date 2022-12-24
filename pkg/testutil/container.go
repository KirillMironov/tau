package testutil

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// MustStartTestContainer starts a test container and terminates it after the test.
func MustStartTestContainer(t *testing.T, request testcontainers.ContainerRequest) testcontainers.Container {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	request.AutoRemove = true

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start test container: %v", err)
	}

	t.Cleanup(func() {
		if err = container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate test container: %v", err)
		}
	})

	return container
}
