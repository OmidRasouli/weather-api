package database

import (
	"context"

	"gorm.io/gorm"
)

// Exec executes a raw SQL query with the provided arguments and returns the result.
func (db *PostgresDB) Exec(query string, args ...interface{}) *gorm.DB {
	return db.conn.Exec(query, args...)
}

// First retrieves the first record that matches the given conditions and maps it to the destination.
func (db *PostgresDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return db.conn.First(dest, conds...)
}

// Create inserts a new record into the database.
func (db *PostgresDB) Create(value interface{}) *gorm.DB {
	return db.conn.Create(value)
}

// Save updates an existing record or inserts a new one if it doesn't exist.
func (db *PostgresDB) Save(value interface{}) *gorm.DB {
	return db.conn.Save(value)
}

// Where adds conditions to a query and returns the result.
func (db *PostgresDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return db.conn.Where(query, args...)
}

// Delete removes a record from the database based on the given conditions.
func (db *PostgresDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	return db.conn.Delete(value, where...)
}

// Table specifies the table to perform operations on.
func (db *PostgresDB) Table(name string) *gorm.DB {
	return db.conn.Table(name)
}

// Preload preloads related data for the specified field.
func (db *PostgresDB) Preload(name string) *gorm.DB {
	return db.conn.Preload(name)
}

// Select specifies the fields to retrieve in a query.
func (db *PostgresDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return db.conn.Select(query, args...)
}

// Updates performs a batch update on records.
func (db *PostgresDB) Updates(values interface{}) *gorm.DB {
	return db.conn.Updates(values)
}

// Find retrieves all records that match the given conditions and maps them to the output.
func (db *PostgresDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	return db.conn.Find(out, where...)
}

// Begin starts a new database transaction.
func (db *PostgresDB) Begin() *gorm.DB {
	return db.conn.Begin()
}

// Commit commits the current transaction.
func (db *PostgresDB) Commit() *gorm.DB {
	return db.conn.Commit()
}

// Rollback rolls back the current transaction.
func (db *PostgresDB) Rollback() *gorm.DB {
	return db.conn.Rollback()
}

// WithContext sets a context for the database operations.
func (db *PostgresDB) WithContext(ctx context.Context) *gorm.DB {
	return db.conn.WithContext(ctx)
}
