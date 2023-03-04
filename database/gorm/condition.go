package gorm

import (
	"fmt"
	"github.com/zedisdog/ty/database"
	tyStrings "github.com/zedisdog/ty/strings"
	"reflect"
	"strings"

	"github.com/zedisdog/ty/errx"
	"gorm.io/gorm"
)

type Condition database.Condition

func (c Condition) Apply(db interface{}) (newDB interface{}, err error) {
	if len(c) < 2 {
		return db, errx.New("condition require at least 2 params")
	}
	query := db.(*gorm.DB)

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
		query = query.Where(s, c[1:]...)
	} else {
		switch len(c) {
		case 2:
			v := reflect.ValueOf(c[1])
			if v.Kind() == reflect.Slice {
				query = query.Where(fmt.Sprintf("%s IN ?", c[0]), c[1])
			} else {
				query = query.Where(fmt.Sprintf("%s = ?", c[0]), c[1])
			}
		case 3:
			query = query.Where(fmt.Sprintf("%s %s (?)", c[0], c[1]), c[2])
		default:
			err = errx.New("condition params is too many")
			return
		}
	}

	newDB = query
	return
}
