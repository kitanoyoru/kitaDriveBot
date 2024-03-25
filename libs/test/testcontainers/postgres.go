package testcontainers

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreatePostgresServerContainer(t *testing.T) (sqlConnectionString string, final func(), err error) {
	_, host, ports, final, err := createContainer(t, testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:latest-v23.1",
		ExposedPorts: []string{"26257/tcp"},
		Cmd:          []string{"start-single-node", "--insecure", "--store=type=mem,size=0.3", "--advertise-addr=localhost", "--log={file-defaults: {dir: /dblogs}, sinks: {stderr: {filter: NONE}}}"},
		WaitingFor: wait.ForSQL("26257/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://root@%s:%s/postgres?sslmode=disable", host, port.Port())
		}).WithStartupTimeout(5 * time.Minute),
	})
	if err != nil {
		return sqlConnectionString, final, err
	}

	sqlConnectionString = fmt.Sprintf("postgres://root@%s:%s/postgres?sslmode=disable", host, ports["26257/tcp"][0].HostPort)
	return sqlConnectionString, final, nil
}
