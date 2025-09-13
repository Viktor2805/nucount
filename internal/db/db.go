package db

import (
	"context"
	"fmt"
	"golang/internal/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
	logger *zap.Logger
}

// NewDatabase initializes a new Database instance.
func New(cfg *config.DatabaseConfig, logger *zap.Logger) (*Db, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: databaseURL,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	d := &Db{DB: db, logger: logger}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping %w", err)
	}

	logger.Info("Database connection established successfully")
	return d, nil
}

// pingDatabase ensures the database connection is alive by pinging it.
func (d *Db) Ping(ctx context.Context) error {
	if d == nil || d.DB == nil {
		return fmt.Errorf("db not initialized")
	}

	sqlDb, err := d.DB.DB()

	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	return sqlDb.PingContext(ctx)
}

// Close disconnects the database connection.
func (d *Db) Close() error {
	if d == nil || d.DB == nil {
		return fmt.Errorf("db not initialized")
	}

	sqlDB, err := d.DB.DB()

	if err != nil {
		return fmt.Errorf("failed to get SQL DB from GORM: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	d.logger.Info("Database connection closed successfully.")
	return nil
}
