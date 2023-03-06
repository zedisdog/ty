package errx

import (
	"errors"
	"fmt"
	tyStrings "github.com/zedisdog/ty/strings"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Error struct {
	code Code
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

// New auto determine the caller and return error with msg.
func New(msg string) error {
	// get call stack, and parse caller by it.
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	arr := strings.Split(string(buf), "\n")
	newArr := arr[:0]
	title := `goroutine \d+ \[running\]:`
	//file := `\t\S+:\d+( \S+|)`
	foot := `\x00+`
	for _, v := range arr {
		if !regexp.MustCompile(title).MatchString(v) &&
			//!regexp.MustCompile(file).MatchString(v) &&
			!regexp.MustCompile(foot).MatchString(v) {
			newArr = append(newArr, v)
		}
	}
	var index int
	for index = 0; index < len(newArr); index += 2 {
		if !tyStrings.ContainersAny(newArr[index], []string{"zedisdog/ty/errx.New", "zedisdog/ty/errx.Wrap", "zedisdog/ty/errx.Make"}) {
			break
		}
	}
	location := strings.Split(strings.Split(strings.Trim(newArr[index+1], "\t"), " ")[0], ":")
	line, err := strconv.Atoi(location[1])
	if err != nil {
		panic(err)
	}

	return &Error{
		msg:  msg,
		file: location[0],
		line: line,
	}
}

func NewWithCode(msg string, code Code) error {
	err := New(msg).(*Error)
	err.code = code
	return err
}
