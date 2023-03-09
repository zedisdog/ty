package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
)

type IMigrator interface {
	Migrate(dsn string, migrates *EmbedDriver) error
}

type DefaultMigrator struct {
}

func (d DefaultMigrator) Migrate(dsn string, migrates *EmbedDriver) (err error) {
	m, err := migrate.NewWithSourceInstance("embedFS", migrates, dsn)
	if err != nil {
		return
	}

	defer func() {
		_, _ = m.Close()
	}()

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}
	return
}
