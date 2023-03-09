package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	urlpkg "net/url"

	"github.com/zedisdog/ty/errx"
)

func buildBody(data interface{}) (body io.ReadCloser, length int64, err error) {
	switch d := data.(type) {
	case []byte:
		body = io.NopCloser(bytes.NewBuffer(d))
		length = int64(len(d))
	case string:
		body = io.NopCloser(bytes.NewBufferString(d))
		length = int64(len(d))
	default:
		if data == nil {
			return
		}
		tmp, err := json.Marshal(data)
		if err != nil {
			return nil, 0, errx.Wrap(err, "covert interface{} to json bytes error")
		}
		body = io.NopCloser(bytes.NewBuffer(tmp))
		length = int64(len(tmp))
	}
	return
}

func WithHeaders(headers map[string][]string) RequestSetter {
	return func(r *http.Request) {
		r.Header = headers
	}
}

type RequestSetter func(*http.Request)

func buildRequest(method string, url string, data interface{}, setters ...RequestSetter) (request *http.Request, err error) {
	u, err := urlpkg.Parse(url)
	if err != nil {
		err = errx.Wrap(err, "parse url error")
		return
	}

	body, length, err := buildBody(data)
	if err != nil {
		err = errx.Wrap(err, "build body error")
		return
	}

	request = &http.Request{
		Method:        method,
		Body:          body,
		URL:           u,
		ContentLength: length,
	}

	for _, setter := range setters {
		setter(request)
	}

	return
}

func PutJSON(url string, data interface{}) (response []byte, err error) {
	return PutWithHeader(url, data, map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	})
}

func PutWithHeader(url string, data interface{}, headers map[string][]string) (response []byte, err error) {
	request, err := buildRequest(http.MethodPut, url, data, WithHeaders(headers))
	if err != nil {
		err = errx.Wrap(err, "put failed")
		return
	}
	return Request(request)
}

// PostJSON post json
//
//	url is the target to post.
//
//	data is to be posted.it can be string, []byte and struct, also nil.
func PostJSON(url string, data interface{}) (response []byte, err error) {
	return PostWithHeader(url, data, map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	})
}

func PostWithHeader(url string, data interface{}, headers map[string][]string) (response []byte, err error) {
	request, err := buildRequest(http.MethodPost, url, data, WithHeaders(headers))
	if err != nil {
		err = errx.Wrap(err, "post failed")
		return
	}
	return Request(request)
}

// GetJSON get json
func GetJSON(url string) (response []byte, err error) {
	return GetWithHeader(url, map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	})
}

func GetWithHeader(url string, headers map[string][]string) (response []byte, err error) {
	request, err := buildRequest(http.MethodGet, url, nil, WithHeaders(headers))
	if err != nil {
		err = errx.Wrap(err, "get failed")
		return
	}

	return Request(request)
}

func Request(request *http.Request) (response []byte, err error) {
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		err = errx.Wrap(err, "request failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		content, _ := io.ReadAll(resp.Body)
		err = errx.NewWithCode(errx.Code(resp.StatusCode), "http error", errx.WithDetail(map[string]string{"content": string(content)}))
		return
	}

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errx.Wrap(err, "read body failed")
		return
	}
	return
}
