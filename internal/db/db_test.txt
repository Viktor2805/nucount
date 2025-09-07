package db_test

import (
	"context"
	"golang/internal/db"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresContainer testcontainers.Container
	postgresUser      = "testuser"
	postgresPassword  = "testpassword"
	postgresDb        = "testdb"
	migrationUrl      = "file://../migrations"
	postgresPort      = "5432"
)

func mustStartPostgresContainer(t *testing.T) error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{postgresPort},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDb,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}
	var err error

	postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("Failed to start Postgres container: %v", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Postgres container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, nat.Port(postgresPort))
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	t.Setenv("POSTGRES_USER", postgresUser)
	t.Setenv("POSTGRES_PASSWORD", postgresPassword)
	t.Setenv("POSTGRES_HOST", host)
	t.Setenv("POSTGRES_PORT", port.Port())
	t.Setenv("POSTGRES_DB", postgresDb)
	t.Setenv("MIGRATION_URL", migrationUrl)

	return nil
}

func teardownTestContainer(t *testing.T) {
	if err := postgresContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("Failed to terminate container: %v", err)
	}
}

func TestDatabase(t *testing.T) {
	err := mustStartPostgresContainer(t)

	if err != nil {
		t.Fatal(err)
	}

	defer teardownTestContainer(t)

	dbInstance, err := db.New()

	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	defer func() {
		sqlDB, err := dbInstance.GetDB().DB()
		if err != nil {
			t.Fatalf("Failed to get SQL DB: %v", err)
		}
		sqlDB.Close()
	}()

	t.Run("Test Database Operations", func(t *testing.T) {
		var result int64
		err := dbInstance.GetDB().Raw("SELECT COUNT(*) FROM information_schema.tables").Scan(&result).Error
		if err != nil {
			t.Fatalf("Failed to execute query: %v", err)
		}
		assert.Greater(t, result, int64(0), "No tables found in the database")
	})
}
