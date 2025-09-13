package db_test

import (
	"context"
	"testing"
	"time"

	"golang/internal/config"
	pkgdb "golang/internal/db"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

const (
	pgUser = "testuser"
	pgPass = "testpassword"
	pgDB   = "testdb"
)

func startPostgres(t *testing.T) (container testcontainers.Container, host string, port int) {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.6",
		Env:          map[string]string{"POSTGRES_USER": pgUser, "POSTGRES_PASSWORD": pgPass, "POSTGRES_DB": pgDB},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort(nat.Port("5432/tcp")),
		).WithDeadline(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	h, err := c.Host(ctx)
	require.NoError(t, err)

	mp, err := c.MappedPort(ctx, nat.Port("5432/tcp"))
	require.NoError(t, err)

	return c, h, mp.Int()
}

func makeDBCfg(host string, port int ) *config.DatabaseConfig {
	return &config.DatabaseConfig{
		User:     pgUser,
		Password: pgPass,
		Host:     host,
		Port:     port,
		Name:     pgDB,
	}
}

func TestNew_OK(t *testing.T) {
	c, host, port := startPostgres(t)
	defer func() { _ = c.Terminate(context.Background()) }()
	cfg := makeDBCfg(host, port)

	logger := zaptest.NewLogger(t)

	db, err := pkgdb.New(cfg, logger)
	require.NoError(t, err, "should connect to containerized Postgres")
	require.NotNil(t, db)

	res := db.Exec("SELECT 1")
	require.NoError(t, res.Error)

	require.NoError(t, db.Close(), "close should succeed")
}

func TestNew_DBDown(t *testing.T) {
	c, host, port := startPostgres(t)
	require.NoError(t, c.Terminate(context.Background()))

	cfg := makeDBCfg(host, port)

	logger := zaptest.NewLogger(t)

	db, err := pkgdb.New(cfg, logger)
	require.Error(t, err, "New must fail because DB is not up")
	require.Nil(t, db)
}

func TestNew_BadPort(t *testing.T) {
	cfg := makeDBCfg("127.0.0.", 1)

	logger := zap.NewNop()

	db, err := pkgdb.New(cfg, logger)
	require.Error(t, err, "should fail to connect with invalid host/port")
	require.Nil(t, db)
}
