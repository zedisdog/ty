package convert

import (
	"database/sql"
	"github.com/jinzhu/copier"
)

func Struct(to interface{}, from interface{}, converters ...copier.TypeConverter) error {
	builtinConverters := []copier.TypeConverter{
		{
			SrcType: sql.NullString{},
			DstType: copier.String,
			Fn:      NullString2String,
		},
		{
			SrcType: copier.String,
			DstType: sql.NullString{},
			Fn:      String2NullString,
		},
	}
	converters = append(converters, builtinConverters...)
	return copier.CopyWithOption(to, from, copier.Option{
		DeepCopy:   true,
		Converters: converters,
	})
}

func NullString2String(src interface{}) (d interface{}, err error) {
	s := src.(sql.NullString)
	if s.Valid {
		d = s.String
	}
	return
}

func String2NullString(src interface{}) (d interface{}, err error) {
	s := src.(string)

	if s == "" {
		d = sql.NullString{String: "", Valid: false}
	} else {
		d = sql.NullString{String: s, Valid: true}
	}

	return
}
