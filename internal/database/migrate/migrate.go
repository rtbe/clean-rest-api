package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// Do runs set of database migrations.
func Do(postgreConn *sqlx.DB) error {
	// Set up database driver for migrations.
	driver, err := postgres.WithInstance(postgreConn.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	// Run all of the migrations from /.migrations folder.
	m.Up()

	return nil
}
