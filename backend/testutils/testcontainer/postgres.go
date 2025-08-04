package testcontainer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"certitrack/internal/config"
	"certitrack/internal/database"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type PostgresContainer struct {
	Container testcontainers.Container
	DB        *gorm.DB
	Config    *config.Config
}

func init() {
	silentLogWriter := &SilentLogger{}
	testcontainers.WithLogger(silentLogWriter)
}

func SetupPostgres(ctx context.Context, cfg *config.Config) (*PostgresContainer, error) {
	testPort := "0" // 0 means system will assign a free port
	if envPort := os.Getenv("POSTGRES_TEST_PORT"); envPort != "" {
		testPort = envPort
	}

	if cfg.Database.User == "" {
		cfg.Database.User = "testuser"
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = "testpassword"
	}
	if cfg.Database.Name == "" {
		cfg.Database.Name = "testdb"
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.Database.User,
			"POSTGRES_PASSWORD": cfg.Database.Password,
			"POSTGRES_DB":       cfg.Database.Name,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(10*time.Second),
		),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = false
			hostConfig.PortBindings = nat.PortMap{
				"5432/tcp": []nat.PortBinding{{
					HostIP:   "0.0.0.0",
					HostPort: testPort,
				}},
			}
		},
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating container: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	cfg.Database.Host = host
	cfg.Database.Port = port.Port()
	cfg.Database.SSLMode = "disable"

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		return nil, err
	}

	return &PostgresContainer{
		Container: postgresContainer,
		DB:        db,
		Config:    cfg,
	}, nil
}

func (pc *PostgresContainer) Teardown(ctx context.Context) error {
	if pc.Container != nil {
		if err := pc.Container.Terminate(ctx); err != nil {
			log.Printf("Error al detener el contenedor: %v", err)
			return err
		}
	}

	// if pc.DB != nil {
	// 	sqlDB, err := pc.DB.DB()
	// 	if err == nil {
	// 		if err := sqlDB.Close(); err != nil {
	// 			log.Printf("Error al cerrar la conexión a la base de datos: %v", err)
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
