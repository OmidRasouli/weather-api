package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB is a wrapper around the GORM database connection.
// It provides methods to interact with the underlying SQL database.
type PostgresDB struct {
	conn *gorm.DB
}

// NewPostgresDB initializes a new PostgresDB instance using the provided DSN (Data Source Name).
// It returns a Database interface implementation or an error if the connection fails.
func NewPostgresDB(dsn string) (interfaces.Database, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &PostgresDB{conn: conn}, nil
}

// DB retrieves the underlying *sql.DB instance from the GORM connection.
// This can be used for lower-level database operations.
func (db *PostgresDB) DB() (*sql.DB, error) {
	return db.conn.DB()
}

// Ping checks the connectivity to the database by sending a ping request.
// It uses a context with a timeout to avoid hanging indefinitely.
func (db *PostgresDB) Ping(ctx context.Context) error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
