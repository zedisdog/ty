package errx

import (
	"errors"
	"fmt"
)

func Wrap(err error, msg string) error {
	return WrapWithCode(err, msg, Nil)
}

func WrapWithCode(err error, msg string, code Code) error {
	if err == nil {
		return nil
	}
	e := New(msg).(*Error)
	e.err = err
	e.code = code
	return e
}

func MakeCodeWrapperWithPrefix(prefix string) func(error, string, Code) error {
	return func(err error, s string, code Code) error {
		return WrapWithCode(err, fmt.Sprintf("[%s]%s", prefix, s), code)
	}
}

// Is reports whether any error in err's tree matches target.
//
// code.Nil is not equal
func Is(err error, target error) bool {
	if errors.Is(err, target) {
		return true
	}

	errxTarget, ok := target.(*Error)
	if !ok || errxTarget.code == Nil {
		return false
	}

	return IsCode(err, errxTarget.code)
}

// IsCode reports whether any error in err's tree matches target code.
func IsCode(err error, target Code) bool {
	if err == nil && target == Nil {
		return false
	}

	_, ok := err.(*Error)
	if !ok {
		return false
	}

	equal := false

	walk(err, func(c Code) bool {
		if c == target {
			equal = true
			return false
		}
		return true
	})

	return equal
}

// walk gets all codes in err tree layer by layer and puts it to handler.
// it stops when handler returns false or there is no inner layer.
func walk(err error, handler func(c Code) bool) {
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
