package application

import (
	"reflect"
	"testing"
)

type isay interface {
	say()
}

type saya struct {
	content string
}

func (sa saya) say() {
	println(sa.content)
}

func TestNormal(t *testing.T) {
	o := interface{}(&saya{content: "a"})
	old := reflect.ValueOf(o).Elem().Interface()

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(interface{}(&saya{content: "b"})).Elem())
	o.(isay).say()

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(old))
	o.(isay).say()
}
