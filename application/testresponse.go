package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"

	"github.com/tidwall/gjson"
)

type TestResponse struct {
	*httptest.ResponseRecorder
}

func (t TestResponse) ParseBody(v interface{}) {
	body, err := io.ReadAll(t.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		panic(err)
	}
}

func (t TestResponse) Body() []byte {
	body, err := io.ReadAll(t.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	t.ResponseRecorder.Body = bytes.NewBuffer(body)
	return body
}

func (t TestResponse) NotEmpty(field string) bool {
	value := gjson.Get(string(t.Body()), field)
	switch value.Type {
	case gjson.False:
		fallthrough
	case gjson.Null:
		return false
	case gjson.Number:
		return value.Num != 0
	case gjson.String:
		return value.String() != ""
	case gjson.JSON:
		return value.Raw != ""
	default:
		panic(errors.New("not a json"))
	}
}

func (t TestResponse) Len(field string) int {
	value := gjson.Get(string(t.Body()), field)
	if !value.IsArray() {
		panic(errors.New("not a json array"))
	}
	return len(value.Array())
}

func (t TestResponse) Get(field string) gjson.Result {
	return gjson.Get(string(t.Body()), field)
}
