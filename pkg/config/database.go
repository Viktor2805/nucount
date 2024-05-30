package config

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	user := Envs.POSTGRES_USER
	password := Envs.POSTGRES_PASSWORD
	host := Envs.POSTGRES_HOST
	port := Envs.POSTGRES_PORT
	dbName := Envs.POSTGRES_DB

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Europe/Kiev",
		user,
		password,
		host,
		port,
		dbName,
	)

	d, err := gorm.Open(postgres.New(postgres.Config{
		DSN: databaseUrl,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := d.DB()
	if err != nil {
		log.Fatalf("failed to get database connection pool: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	DB = d

	RunMigrations(databaseUrl)
}

func RunMigrations(databaseUrl string) {
	m, err := migrate.New(
		Envs.MIGRATION_URL,
		databaseUrl,
	)

	if err != nil {
		log.Fatal("cannot create migrate instance", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("migration is run successfully")
}
