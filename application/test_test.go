package application

import (
	"reflect"
	"sync"
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
	m := new(sync.Map)
	m.Store("say", &saya{content: "a"})

	o, _ := m.Load("say")
	old := reflect.ValueOf(o).Elem().Interface()

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(&saya{content: "b"}).Elem())
	o.(isay).say()

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(old))
	o.(isay).say()
}
