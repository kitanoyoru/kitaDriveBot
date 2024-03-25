package testcontainers

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

func createContainer(t *testing.T, req testcontainers.ContainerRequest) (container testcontainers.Container, host string, ports nat.PortMap, final func(), err error) {
	t.Helper()

	final = func() {}
	ctx := context.Background()

	container, err = testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return container, host, ports, final, err
	}

	host, err = container.Host(ctx)
	if err != nil {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
		return container, host, ports, final, err
	}

	ports, err = container.Ports(ctx)
	if err != nil {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
		return container, host, ports, final, err
	}

	final = func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return container, host, ports, final, nil
}
