//go:build sqlite

package gorm

import (
	"github.com/zedisdog/ty/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(path string) (DB database.IDatabase, err error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		err = database.Wrap(err, "connect database using gorm failed")
		return
	}
	DB = &Database{
		db: db,
	}
	return
}
