package gorm

import "github.com/zedisdog/ty/errx"

func Wrap(err error, msg string) error {
	return WrapWithCode(err, msg, errx.Nil)
}

func WrapWithCode(err error, msg string, code errx.Code) error {
	return errx.WrapWithCode(err, "database: "+msg, code)
}
