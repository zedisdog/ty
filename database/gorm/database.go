package gorm

import (
	"github.com/zedisdog/ty/database"
	"github.com/zedisdog/ty/errx"
	"gorm.io/gorm"
	"reflect"
)

type Database struct {
	db *gorm.DB
}

// DB get the gorm db instance
func (d *Database) DB() *gorm.DB {
	return d.db
}

func (d *Database) Create(model interface{}) error {
	return Wrap(d.db.Create(model).Error, "create failed")
}

func (d Database) Where(conditions ...database.Condition) database.IDatabase {
	var (
		db  interface{}
		err error
	)
	for _, condition := range conditions {
		db, err = (Condition)(condition).Apply(d.db)
		if err != nil {
			panic(Wrap(err, "condition apply failed"))
		}
	}
	d.db = db.(*gorm.DB)
	return &d
}

func (d *Database) Update(model interface{}, m map[string]interface{}) (count int64, err error) {
	result := d.db.Model(model).Updates(m)
	return result.RowsAffected, Wrap(result.Error, "update failed")
}

func (d *Database) Delete(model interface{}) error {
	return Wrap(d.db.Delete(model).Error, "delete failed")
}

func (d *Database) First(model interface{}) (err error) {
	err = d.db.First(&model).Error
	if err == gorm.ErrRecordNotFound {
		err = WrapWithCode(err, "not found", errx.NotFound)
	} else {
		err = Wrap(err, "find first record failed")
	}
	return
}

func (d *Database) Find(list interface{}) (err error) {
	return Wrap(d.db.Find(&list).Error, "get list failed")
}

func (d *Database) Page(page int, size int, list interface{}) (total int64, err error) {
	t := reflect.TypeOf(list).Elem().Elem()
	err = Wrap(d.db.Model(reflect.New(t)).Count(&total).Error, "get count failed")
	if err != nil {
		return
	}

	err = Wrap(d.db.Offset((page-1)*size).Limit(size).Find(list).Error, "get page list failed")
	return
}

func (d *Database) Transaction(f func(tx database.IDatabase) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return f(&Database{
			db: tx,
		})
	})
}
