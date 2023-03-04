package errx

import (
	"errors"
	"fmt"
	"github.com/zedisdog/ty/errx/code"
	"runtime"
	"strings"
)

type Error struct {
	code code.ICode
	msg  string
	err  error
	file string
	line int
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e.file, e.line, e.msg)))
		case s.Flag('#'):
			strArr := []string{fmt.Sprintf("%+v", e)}

			err := errors.Unwrap(e)
			for {
				if err == nil {
					break
				}

				if _, ok := err.(*Error); ok {
					strArr = append(strArr, fmt.Sprintf("%+v", err))
				} else {
					strArr = append(strArr, err.Error()+"\n")
				}

				err = errors.Unwrap(err)
			}
			_, _ = s.Write([]byte(strings.Join(strArr, "")))
		default:
			_, _ = s.Write([]byte(e.Error()))
		}
	}
}

func New(msg string) error {
	return NewSkip(msg, 2)
}

func NewWithCode(msg string, code code.ICode) error {
	err := NewSkip(msg, 2).(*Error)
	err.code = code
	return err
}

func NewSkip(msg string, skip int) error {
	_, file, line, _ := runtime.Caller(skip)
	return &Error{
		msg:  msg,
		file: file,
		line: line,
	}
}
