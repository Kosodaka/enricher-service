package integration

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

// PostgreSQLContainer wraps testcontainers.Container with extra methods.
type PostgreSQLContainer struct {
	testcontainers.Container
	Config PostgreSQLContainerConfig
}

type PostgreSQLContainerOption func(c *PostgreSQLContainerConfig)

type PostgreSQLContainerConfig struct {
	User     string
	Password string
	Port     string
	Database string
	Host     string
}

// GetDSN returns DB connection URL.
func (c PostgreSQLContainer) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Config.User, c.Config.Password, c.Config.Host, c.Config.Port, c.Config.Database)
}

// NewPostgreSQLContainer creates and starts a PostgreSQL container.
func NewPostgreSQLContainer(ctx context.Context, opts ...PostgreSQLContainerOption) (*PostgreSQLContainer, error) {
	const (
		psqlImage = "postgres"
		psqlPort  = "5432"
	)

	// Define container ENVs
	config := PostgreSQLContainerConfig{
		User:     "admin",
		Password: "qwerty",
		Database: "human",
	}
	for _, opt := range opts {
		opt(&config)
	}

	containerPort := psqlPort + "/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"POSTGRES_HOST":     config.Host,
				"POSTGRES_DB":       config.Database,
				"POSTGRES_PORT":     containerPort,
				"POSTGRES_USER":     config.User,
				"POSTGRES_PASSWORD": config.Password,
			},
			ExposedPorts: []string{
				containerPort,
			},

			Image: fmt.Sprintf("%s", psqlImage),
			WaitingFor: wait.ForExec([]string{"pg_isready", "-d", config.Database, "-U", config.User}).
				WithPollInterval(1 * time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
					return exitCode == 0
				}),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting request provider: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("getting port for (%s): %w", containerPort, err)
	}
	config.Port = mappedPort.Port()
	config.Host = host

	fmt.Println("Host:", config.Host, config.Port)

	return &PostgreSQLContainer{
		Container: container,
		Config:    config,
	}, nil
}
