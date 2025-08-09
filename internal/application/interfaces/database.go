package interfaces

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// Database defines an interface for interacting with a database.
// It provides a set of methods for common database operations such as querying, creating, updating, and deleting records.
type Database interface {
	// Exec executes a raw SQL query with the provided arguments and returns a *gorm.DB instance.
	Exec(query string, args ...interface{}) *gorm.DB

	// DB retrieves the underlying *sql.DB instance for lower-level database operations.
	DB() (*sql.DB, error)

	// Ping checks the connectivity to the database using the provided context.
	Ping(ctx context.Context) error

	// Create inserts a new record into the database.
	Create(value interface{}) *gorm.DB

	// First retrieves the first record that matches the given conditions.
	First(out interface{}, where ...interface{}) *gorm.DB

	// Save updates an existing record or inserts a new one if it doesn't exist.
	Save(value interface{}) *gorm.DB

	// Delete removes a record from the database based on the given conditions.
	Delete(value interface{}, where ...interface{}) *gorm.DB

	// Table specifies the table to perform operations on.
	Table(name string) *gorm.DB

	// Preload preloads related data for the specified field.
	Preload(name string) *gorm.DB

	// Select specifies the fields to retrieve in a query.
	Select(query interface{}, args ...interface{}) *gorm.DB

	// Where adds conditions to a query.
	Where(query interface{}, args ...interface{}) *gorm.DB

	// Updates performs a batch update on records.
	Updates(values interface{}) *gorm.DB

	// Find retrieves all records that match the given conditions.
	Find(out interface{}, where ...interface{}) *gorm.DB

	// Begin starts a new database transaction.
	Begin() *gorm.DB

	// Commit commits the current transaction.
	Commit() *gorm.DB

	// Rollback rolls back the current transaction.
	Rollback() *gorm.DB

	// WithContext sets a context for the database operations.
	WithContext(ctx context.Context) *gorm.DB

	// Close closes the database connection.
	Close() error
}
