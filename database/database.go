package database

import (
	"fmt"
	"github.com/zedisdog/ty/errx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

func NewDatabase(dsn string) (db *gorm.DB, err error) {
	var (
		config = strings.Split(dsn, "://")
	)

	switch config[0] {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config[1]), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
					LogLevel:                  logger.Warn,            // Log level
					IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
					Colorful:                  false,                  // Disable color
				},
			),
		})
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
