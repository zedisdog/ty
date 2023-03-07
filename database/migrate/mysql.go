package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"io/fs"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/zedisdog/ty/errx"
)

func InitAutoMigrateForMysqlFunc(dsn string, fss ...fs.FS) func() error {
	if len(fss) < 1 {
		panic(errx.New("fs can not be empty"))
	}
	return func() (err error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return
		}
		defer db.Close()

		instance, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return
		}

		driver := NewFsDriver()
		for _, f := range fss {
			driver.Add(f)
		}

		m, err := migrate.NewWithInstance("", driver, "main", instance)
		if err != nil {
			return
		}

		err = m.Up()
		if err != nil && err == migrate.ErrNoChange {
			err = nil
		}
		return
	}
}
