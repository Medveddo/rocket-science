package postgres

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

const (
	DefaultPostgresImage    = "postgres:16-alpine"
	DefaultPostgresPort     = "5432"
	DefaultPostgresUser     = "postgres"
	DefaultPostgresPassword = "postgres"
	DefaultPostgresDatabase = "test"
)

type Container struct {
	container testcontainers.Container
	config    *Config
	logger    Logger
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	ImageName     string
	ContainerName string
	Port          string
	User          string
	Password      string
	Database      string
	NetworkName   string
	Logger        Logger
	Host          string
	WaitStrategy  wait.Strategy
}

type Option func(*Config)

func WithImageName(imageName string) Option {
	return func(c *Config) {
		c.ImageName = imageName
	}
}

func WithContainerName(name string) Option {
	return func(c *Config) {
		c.ContainerName = name
	}
}

func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithUser(user string) Option {
	return func(c *Config) {
		c.User = user
	}
}

func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

func WithNetworkName(networkName string) Option {
	return func(c *Config) {
		c.NetworkName = networkName
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithWaitStrategy(waitStrategy wait.Strategy) Option {
	return func(c *Config) {
		c.WaitStrategy = waitStrategy
	}
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	config := &Config{
		ImageName:     DefaultPostgresImage,
		ContainerName: "test-postgres",
		Port:          DefaultPostgresPort,
		User:          DefaultPostgresUser,
		Password:      DefaultPostgresPassword,
		Database:      DefaultPostgresDatabase,
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(config)
	}

	req := testcontainers.ContainerRequest{
		Image:        config.ImageName,
		Name:         config.ContainerName,
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", config.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Database,
		},
		Networks:   []string{config.NetworkName},
		WaitingFor: config.WaitStrategy,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	config.Host = host
	config.Port = port.Port()

	config.Logger.Info(ctx, "PostgreSQL container started successfully",
		zap.String("container", config.ContainerName),
		zap.String("image", config.ImageName),
		zap.String("host", config.Host),
		zap.String("port", config.Port),
	)

	return &Container{
		container: container,
		config:    config,
		logger:    config.Logger,
	}, nil
}

func (c *Container) Container() testcontainers.Container {
	return c.container
}

func (c *Container) Config() *Config {
	return c.config
}

func (c *Container) Host(ctx context.Context) (string, error) {
	return c.config.Host, nil
}

func (c *Container) Port(ctx context.Context) (string, error) {
	return c.config.Port, nil
}

func (c *Container) ConnectionString(ctx context.Context) (string, error) {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.config.User, c.config.Password, c.config.Host, c.config.Port, c.config.Database), nil
}

func (c *Container) Stop(ctx context.Context) error {
	c.logger.Info(ctx, "Stopping PostgreSQL container", zap.String("container", c.config.ContainerName))
	return c.container.Terminate(ctx)
}
