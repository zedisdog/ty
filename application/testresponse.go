package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http/httptest"
)

type TestResponse struct {
	*httptest.ResponseRecorder
}

func (t TestResponse) ParseBody(v interface{}) {
	err := json.Unmarshal(t.Body(), v)
	if err != nil {
		panic(err)
	}
}

func (t TestResponse) ParseField(field string, v interface{}) {
	switch vv := v.(type) {
	case *string:
		*vv = t.Get(field).String()
	case *int64:
		*vv = t.Get(field).Int()
	default:
		err := json.Unmarshal([]byte(t.Get(field).String()), v)
		if err != nil {
			panic(err)
		}
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
