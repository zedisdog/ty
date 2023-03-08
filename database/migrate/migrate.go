package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/zedisdog/ty/errx"
)

type IMigrator interface {
	Migrate(dbType string, db *sql.DB, migrates *EmbedDriver) error
}

type DefaultMigrator struct {
}

func (d DefaultMigrator) Migrate(dbType string, db *sql.DB, migrates *EmbedDriver) (err error) {
	var (
		driver database.Driver
	)
	switch dbType {
	case "mysql":
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	case "sqlite3":
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	default:
		err = errx.New("unsupported database type")
	}
	if err != nil {
		return
	}

	m, err := migrate.NewWithInstance("database", migrates, "dbType", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}
	return
}
