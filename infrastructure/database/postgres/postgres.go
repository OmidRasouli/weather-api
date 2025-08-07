package postgres

import (
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/pkg/logger"
)

// PostgresConfig holds the configuration details required to connect to a PostgreSQL database.
type PostgresConfig struct {
	// Host specifies the database server address.
	Host string
	// Port specifies the port on which the database server is listening.
	Port string
	// User specifies the username for database authentication.
	User string
	// Password specifies the password for database authentication.
	Password string
	// DBName specifies the name of the database to connect to.
	DBName string
	// SSLMode specifies the SSL mode for the connection (e.g., "disable", "require").
	SSLMode string
}

// GetDSN constructs the Data Source Name (DSN) string from the PostgresConfig fields.
func (pc PostgresConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.Host,
		5432,
		pc.User,
		pc.Password,
		pc.DBName,
		pc.SSLMode,
	)
}

// NewPostgresConnection establishes a connection to the PostgreSQL database using the provided configuration.
// It retries the connection up to a maximum number of attempts in case of failure.
func NewPostgresConnection(config PostgresConfig) (database.Database, error) {
	// Generate the DSN string from the configuration.
	dsn := config.GetDSN()

	// Initialize variables for the database instance and error handling.
	var dbInstance database.Database
	var err error
	maxRetries := 5

	// Attempt to connect to the database with retries.
	for i := 0; i < maxRetries; i++ {
		// Try to create a new database connection.
		dbInstance, err = database.NewPostgresDB(dsn)
		if err == nil {
			break
		}

		// Log the error and wait before retrying.
		logger.Errorf("Failed to connect to database. Retry %d/%d after 5 seconds...", i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}

	// If all retries fail, return an error.
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, err)
	}

	// Configure the database connection pool settings.
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set the maximum number of idle connections.
	sqlDB.SetMaxIdleConns(10)
	// Set the maximum number of open connections.
	sqlDB.SetMaxOpenConns(100)
	// Set the maximum lifetime of a connection.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Return the established database instance.
	return dbInstance, nil
}
