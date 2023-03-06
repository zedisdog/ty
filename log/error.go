package log

import "github.com/zedisdog/ty/errx"

func Wrap(err error, msg string) error {
	return WrapWithCode(err, msg, errx.Nil)
}

var WrapWithCode = errx.MakeCodeWrapperWithPrefix("log")
