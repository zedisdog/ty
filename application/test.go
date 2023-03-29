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
	IHasDatabase
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

func TestSuit(f func(h ICanTest)) {
	GetInstance().TestSuit(f)
}
func (app *App) TestSuit(f func(h ICanTest)) {
	app.config.Set("mode", "testing")
	//数据库的引用，比如 *gorm.DB
	pDB := app.Database("default")

	//拷贝了一份值
	old := reflect.ValueOf(pDB).Elem().Interface()

	//调用Begin方法, 如果没有这个方法说明有问题, 它自己会panic
	tx := reflect.ValueOf(pDB).MethodByName("Begin").Call(nil)[0]

	//把pDB引用指向的值替换成tx
	reflect.ValueOf(pDB).Elem().Set(tx.Elem())

	f(app)

	//tx.Rollback
	tx.MethodByName("Rollback").Call(nil)

	//再把pDB引用指向的值替换回来, 无事发生
	reflect.ValueOf(pDB).Elem().Set(reflect.ValueOf(old))
}
