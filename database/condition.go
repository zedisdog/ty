package database

import (
	"fmt"
	tyStrings "github.com/zedisdog/ty/strings"
	"reflect"
	"strings"

	"github.com/zedisdog/ty/errx"
	"gorm.io/gorm"
)

// Condition present database query condition
//
//	e.g:
//		Condition{"field", "123"}:                 field = "123"
//		Condition{"field", ">", 100}:              field > 100
//		Condition{"filed1 = 100 OR field2 = 200"}: filed1 = 100 OR field2 = 200
type Condition []interface{}

func (c Condition) MustApply(db *gorm.DB) *gorm.DB {
	newDB, err := c.Apply(db)
	if err != nil {
		panic(err)
	}
	return newDB
}

func (c Condition) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	if len(c) < 2 {
		return db, errx.New("condition require at least 2 params")
	}
	newDB = db

	list := []string{
		" and ",
		" or ",
		"?",
		" not ",
		" between ",
		" like ",
		" is ",
	}
	if s, ok := c[0].(string); ok && tyStrings.ContainersAny(strings.ToLower(s), list) {
		newDB = newDB.Where(s, c[1:]...)
	} else {
		switch len(c) {
		case 2:
			v := reflect.ValueOf(c[1])
			if v.Kind() == reflect.Slice {
				newDB = newDB.Where(fmt.Sprintf("%s IN ?", c[0]), c[1])
			} else {
				newDB = newDB.Where(fmt.Sprintf("%s = ?", c[0]), c[1])
			}
		case 3:
			newDB = newDB.Where(fmt.Sprintf("%s %s (?)", c[0], c[1]), c[2])
		default:
			err = errx.New("condition params is too many")
			return
		}
	}

	return
}

type Conditions []Condition

func (cs Conditions) MustApply(db *gorm.DB) *gorm.DB {
	newDB, err := cs.Apply(db)
	if err != nil {
		panic(err)
	}

	return newDB
}

func (cs Conditions) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	newDB = db
	for _, c := range cs {
		newDB, err = c.Apply(newDB)
		if err != nil {
			return
		}
	}
	return
}
