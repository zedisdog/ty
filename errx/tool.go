package errx

import (
	"errors"
	"github.com/zedisdog/ty/errx/code"
)

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	e := NewSkip(msg, 2).(*Error)
	e.err = err
	return e
}

// Is reports whether any error in err's tree matches target.
//
// code.Nil is not equal
func Is(err error, target error) bool {
	if errors.Is(err, target) {
		return true
	}

	errxTarget, ok := target.(*Error)
	if !ok || errxTarget.code == code.Nil {
		return false
	}

	equal := false
	walk(err, func(c code.ICode) bool {
		if c == errxTarget.code {
			equal = true
			return false
		}
		return true
	})

	return equal
}

// IsCode reports whether any error in err's tree matches target code.
func IsCode(err error, target code.ICode) bool {
	if err == nil {
		return target == code.Nil
	}

	_, ok := err.(*Error)
	if !ok {
		return target == code.Nil
	}

	equal := false

	walk(err, func(c code.ICode) bool {
		if c == target {
			equal = true
			return false
		}
		return true
	})

	return equal
}

func walk(err error, handler func(c code.ICode) bool) {
	e := err
	for e != nil {
		errxErr, ok := e.(*Error)
		if ok && !handler(errxErr.code) {
			break
		}

		if ee, ok := e.(interface{ Unwrap() error }); ok {
			e = ee.Unwrap()
		} else {
			e = nil
		}
	}
}
