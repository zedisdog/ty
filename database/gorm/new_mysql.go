//go:build !sqlite

package gorm

import (
	"github.com/zedisdog/ty/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) (DB database.IDatabase, err error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		err = database.Wrap(err, "connect database using gorm failed")
		return
	}
	DB = &Database{
		db: db,
	}
	return
}
