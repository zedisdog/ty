package errx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err := New("test").(*Error)

	require.Equal(t, file, err.file)
	require.Equal(t, line+2, err.line)
}

func TestWrap(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err := New("test1")
	err2 := Wrap(err, "test2").(*Error)

	require.Equal(t, file, err2.file)
	require.Equal(t, line+3, err2.line)
}

func TestFormat(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err1 := New("test1")
	err2 := Wrap(err1, "test2")
	err3 := Wrap(err2, "test3")

	require.Equal(t, "test3", fmt.Sprintf("%v", err3))
	require.Equal(t, fmt.Sprintf("%s:%d:test3\n", file, line+4), fmt.Sprintf("%+v", err3))

	except := []string{
		fmt.Sprintf("%s:%d:%s\n", file, line+4, err3.Error()),
		fmt.Sprintf("%s:%d:%s\n", file, line+3, err2.Error()),
		fmt.Sprintf("%s:%d:%s\n", file, line+2, err1.Error()),
	}
	require.Equal(t, strings.Join(except, ""), fmt.Sprintf("%#v", err3))
}
