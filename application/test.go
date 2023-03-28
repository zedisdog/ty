package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zedisdog/ty/auth"
	"github.com/zedisdog/ty/errx"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type ICanTest interface {
	Request(method string, path string, header http.Header, body io.Reader) (w *TestResponse, req *http.Request)
	Get(path string) (*TestResponse, *http.Request)
	Delete(path string) (*TestResponse, *http.Request)
	Post(path string, data any) (*TestResponse, *http.Request)
	Put(path string, data any) (*TestResponse, *http.Request)
	WithHeader(header http.Header) ICanTest
	ActingAs(id uint64, claims ...map[string]interface{}) ICanTest
	TestSuit(f func(h ICanTest))
}

func (app App) WithHeader(header http.Header) ICanTest {
	app.header = header
	return &app
}

func (app App) Request(method string, path string, header http.Header, body io.Reader) (w *TestResponse, req *http.Request) {
	w = &TestResponse{ResponseRecorder: httptest.NewRecorder()}
	req = httptest.NewRequest(method, path, body)
	req.Header = app.mergeHeader(
		map[string][]string{
			"Content-Type": {"application/json"},
		},
		header,
	)
	svc, ok := app.httpServers.Load("default")
	if !ok {
		panic(errx.New("theres no server"))
	}
	svc.(http.Handler).ServeHTTP(w.ResponseRecorder, req)
	return
}

func (app App) Get(path string) (*TestResponse, *http.Request) {
	return app.Request(http.MethodGet, path, app.header, nil)
}

func (app App) Delete(path string) (*TestResponse, *http.Request) {
	return app.Request(http.MethodDelete, path, app.header, nil)
}

func (app App) Post(path string, data any) (*TestResponse, *http.Request) {
	body, err := app.buildBody(data)
	if err != nil {
		panic(err)
	}

	return app.Request(http.MethodPost, path, app.header, body)
}

func (app App) Put(path string, data any) (*TestResponse, *http.Request) {
	body, err := app.buildBody(data)
	if err != nil {
		panic(err)
	}

	return app.Request(http.MethodPut, path, app.header, body)
}

func (app App) buildBody(data any) (body io.Reader, err error) {
	if data == nil {
		return
	}

	var d []byte
	switch dd := data.(type) {
	case string:
		d = []byte(dd)
	default:
		d, err = json.Marshal(dd)
		if err != nil {
			return
		}
	}
	body = bytes.NewBuffer(d)
	return
}

func (app App) mergeHeader(origin http.Header, delta http.Header) http.Header {
	if origin == nil {
		origin = make(http.Header)
	}
	if delta == nil {
		return origin
	}
	for key, value := range delta {
		origin[key] = value
	}

	return origin
}

func (app App) ActingAs(id uint64, claims ...map[string]interface{}) ICanTest {
	builder := auth.NewJwtTokenBuilder().WithKey(app.config.GetString("http.authKey")).WithClaim(auth.JwtSubject, gconv.String(id))
	for _, c := range claims {
		builder = builder.WithClaims(c)
	}
	token, err := builder.BuildToken()
	if err != nil {
		panic(err)
	}
	return app.WithHeader(map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}})
}

type IHasTransaction interface {
	Begin() IHasTransaction
	Rollback()
}

func (app *App) TestSuit(f func(h ICanTest)) {
	app.config.Set("mode", "testing")
	old := reflect.ValueOf(app.Database("default")).Elem().Interface()

	reflect.ValueOf(app.Database("default")).Elem().
		Set(reflect.ValueOf(old.(IHasTransaction).Begin()).Elem())

	f(app)

	app.Database("default").(IHasTransaction).Rollback()
	reflect.ValueOf(app.Database("default")).Elem().
		Set(reflect.ValueOf(old))
}
