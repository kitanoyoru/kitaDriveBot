package testcontainers

import (
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateMinioServerContainer(t *testing.T) (minioHost string, final func(), err error) {
	_, host, ports, final, err := createContainer(t, testcontainers.ContainerRequest{
		Image:        "minio/minio:RELEASE.2023-07-21T21-12-44Z",
		ExposedPorts: []string{"9000/tcp"},
		Cmd:          []string{"server", "/data"},
		WaitingFor:   wait.ForListeningPort("9000/tcp").WithStartupTimeout(5 * time.Minute),
	})
	if err != nil {
		return minioHost, final, err
	}

	minioHost = fmt.Sprintf("%s:%s", host, ports["9000/tcp"][0].HostPort)

	return minioHost, final, nil
}
