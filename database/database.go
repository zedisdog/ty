package database

import (
	"database/sql"
	"github.com/zedisdog/ty/errx"
)

const NotFound = 404

func init() {
	err := errx.Register(NotFound, "not found")
	if err != nil {
		panic(err)
	}
}

// Condition present database query condition
//
//	e.g:
//		Condition{"field", "123"}:                 field = "123"
//		Condition{"field", ">", 100}:              field > 100
//		Condition{"filed1 = 100 OR field2 = 200"}: filed1 = 100 OR field2 = 200
type Condition []interface{}

type IDatabase interface {
	Create(interface{}) error
	// Where Copy the instance and set query conditions
	Where(conditions ...Condition) IDatabase
	Update(interface{}, map[string]interface{}) (count int64, err error)
	UpdateModel(model interface{}) (count int64, err error)
	Delete(interface{}) error
	// First finds one record
	First(interface{}) error
	// Find finds multi record
	Find(interface{}) error
	// Page Finds multi record with pagination
	Page(page int, size int, list interface{}) (total int64, err error)
	Transaction(f func(IDatabase) error) error
	RawDB() (*sql.DB, error)
	Exec(sql string, args ...interface{}) IDatabase
	Scan(model interface{}) error
}
