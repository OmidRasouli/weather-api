package database

import (
	"database/sql"
	"fmt"

	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

// TODO
// migration command
// migrate create -ext sql -dir db/migrations -seq create_users_table

type DBInstance interface {
	DB() (*sql.DB, error)
}

type Migration interface {
	Up() error
	Steps(int) error
	Migrate(version uint) error
	Version() (uint, bool, error)
}

type MigrationManager struct {
	db         DBInstance
	migrations Migration
}

func NewMigrationManager(db DBInstance, migrations Migration) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: migrations,
	}
}

func NewMigrateInstance(db DBInstance, migrationsPath string, dbname string) (Migration, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance("file://"+migrationsPath, dbname, driver)
}

// RunMigrations applies all pending migrations to the database.
func (mm *MigrationManager) RunMigrations() error {
	// Get current version before running migrations
	currentVersion, dirty, err := mm.migrations.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	// Handle dirty state
	if dirty {
		logger.Warnf("WARNING: Database is in dirty state at version %d", currentVersion)
		if err := mm.handleDirtyState(); err != nil {
			return fmt.Errorf("failed to handle dirty state: %w", err)
		}
	}

	// Apply migrations
	if err := mm.migrations.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("Database is up to date")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get new version after migrations
	newVersion, _, err := mm.migrations.Version()
	if err != nil {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	logger.Infof("Successfully migrated from version %d to %d", currentVersion, newVersion)
	return nil
}

// handleDirtyState handles the scenario where the database is in a dirty state.
// This method should include logic to resolve the dirty state, such as taking backups or forcing a version.
func (mm *MigrationManager) handleDirtyState() error {
	// Implement your dirty state handling logic here
	// For example, you might want to:
	// 1. Take a backup
	// 2. Force the version
	// 3. Notify administrators
	return nil
}

// RollbackLastMigration rolls back the last applied migration.
func (mm *MigrationManager) RollbackLastMigration() error {
	if err := mm.migrations.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback last migration: %w", err)
	}
	logger.Info("Successfully rolled back last migration")
	return nil
}

// RollbackToVersion rolls back the database to a specific migration version.
func (mm *MigrationManager) RollbackToVersion(version uint) error {
	if err := mm.migrations.Migrate(version); err != nil {
		return fmt.Errorf("failed to rollback to version %d: %w", version, err)
	}
	logger.Infof("Successfully rolled back to version %d", version)
	return nil
}

// GetMigrationStatus retrieves the current migration version and dirty state of the database.
func (mm *MigrationManager) GetMigrationStatus() (uint, bool, error) {
	version, dirty, err := mm.migrations.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration status: %w", err)
	}
	return version, dirty, nil
}
