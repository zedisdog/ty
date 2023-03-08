package gorm

import (
	"database/sql"
	"github.com/zedisdog/ty/database"
	"github.com/zedisdog/ty/errx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func NewDatabase(dsn string) (DB database.IDatabase, err error) {
	var (
		config []string = strings.Split(dsn, "://")
		db     *gorm.DB
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
		err = database.Wrap(err, "connect database using gorm failed")
		return
	}

	DB = &Database{
		db: db,
	}
	return
}

type Database struct {
	db *gorm.DB
}

// DB get the gorm db instance
func (d *Database) DB() *gorm.DB {
	return d.db
}

func (d *Database) RawDB() (*sql.DB, error) {
	return d.db.DB()
}

func (d *Database) Create(model interface{}) error {
	return database.Wrap(d.db.Create(model).Error, "create failed")
}

func (d Database) Where(conditions ...database.Condition) database.IDatabase {
	var (
		db  interface{}
		err error
	)
	for _, condition := range conditions {
		db, err = (Condition)(condition).Apply(d.db)
		if err != nil {
			panic(database.Wrap(err, "condition apply failed"))
		}
	}
	d.db = db.(*gorm.DB)
	return &d
}

func (d *Database) Update(model interface{}, m map[string]interface{}) (count int64, err error) {
	result := d.db.Model(model).Updates(m)
	return result.RowsAffected, database.Wrap(result.Error, "update failed")
}

func (d *Database) Delete(model interface{}) error {
	return database.Wrap(d.db.Delete(model).Error, "delete failed")
}

func (d *Database) First(model interface{}) (err error) {
	err = d.db.First(&model).Error
	if err == gorm.ErrRecordNotFound {
		err = database.WrapWithCode(err, "not found", database.NotFound)
	} else {
		err = database.Wrap(err, "find first record failed")
	}
	return
}

func (d *Database) Find(list interface{}) (err error) {
	return database.Wrap(d.db.Find(&list).Error, "get list failed")
}

func (d *Database) Page(page int, size int, list interface{}) (total int64, err error) {
	t := reflect.TypeOf(list).Elem().Elem()
	err = database.Wrap(d.db.Model(reflect.New(t)).Count(&total).Error, "get count failed")
	if err != nil {
		return
	}

	err = database.Wrap(d.db.Offset((page-1)*size).Limit(size).Find(list).Error, "get page list failed")
	return
}

func (d *Database) Transaction(f func(tx database.IDatabase) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return f(&Database{
			db: tx,
		})
	})
}
