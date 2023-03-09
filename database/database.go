package database

import (
	"fmt"
	"github.com/zedisdog/ty/errx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

func NewDatabase(dsn string) (db *gorm.DB, err error) {
	var (
		config = strings.Split(dsn, "://")
	)

	switch config[0] {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config[1]), &gorm.Config{})
	case "sqlite":
		fallthrough
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(config[1]), &gorm.Config{})
	default:
		err = errx.New("unsupported database type")
	}

	if err != nil {
		err = errx.Wrap(err, fmt.Sprintf("[database] connect database using gorm failed"))
		return
	}

	return
}
