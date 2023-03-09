package oauth

import (
	"fmt"
	"strings"
)

func WithScope(scope string) func(*redirectOptions) {
	return func(ro *redirectOptions) {
		ro.scope = scope
	}
}

func WithState(state map[string]string) func(*redirectOptions) {
	return func(ro *redirectOptions) {
		ro.state = state
	}
}

type redirectOptions struct {
	scope string
	state map[string]string
}

func (r redirectOptions) State() string {
	tmp := make([]string, 0, len(r.state))
	for key, value := range r.state {
		tmp = append(tmp, fmt.Sprintf("%s=%s", strings.Trim(key, " "), strings.Trim(value, " ")))
	}
	return strings.Join(tmp, "|")
}

type State string

func (s State) Get(key string) (val string) {
	st := strings.Split(string(s), "|")
	for _, item := range st {
		keyAndVal := strings.Split(item, "=")
		if len(keyAndVal) != 2 {
			return
		}
		if keyAndVal[0] == strings.Trim(key, " ") {
			val = keyAndVal[1]
			return
		}
	}
	return
}
