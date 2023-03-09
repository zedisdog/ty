package qrcode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/wechat/mini/auth"
	"github.com/zedisdog/ty/sdk/wechat/mini/qrcode/response"
)

const (
	qrCodeUnlimited = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

type Env string

const (
	Release Env = "release"
	Trial   Env = "trial"
	Develop Env = "develop"
)

func NewQrCode(a *auth.Auth) *QrCode {
	return &QrCode{
		auth: a,
	}
}

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type QrCodeUnlimitedOptions struct {
	Scene      string `json:"scene"`
	Page       Env    `json:"page"`
	CheckPath  bool   `json:"check_path"`
	EnvVersion string `json:"env_version"`
	Width      int    `json:"width"`
	AutoColor  bool   `json:"auto_color"`
	LineColor  Color  `json:"line_color"`
	IsHyaline  bool   `json:"is_hyaline"`
}

func WithCheckPath(enabled bool) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.CheckPath = enabled
	}
}

func WithPage(path Env) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.Page = path
	}
}

func WithWidth(width int) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.Width = width
	}
}

func WithAutoColor(enabled bool) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.AutoColor = enabled
	}
}

func WithLineColor(color Color) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.LineColor = color
	}
}

func WithIsHyaline(enabled bool) func(options *QrCodeUnlimitedOptions) {
	return func(options *QrCodeUnlimitedOptions) {
		options.IsHyaline = enabled
	}
}

func WithEnvVersion(env string) {

}

type QrCode struct {
	auth *auth.Auth
}

func (q QrCode) GetUnlimited(scene map[string]string, setters ...func(*QrCodeUnlimitedOptions)) (r response.QrCodeUnlimited, err error) {
	options := QrCodeUnlimitedOptions{
		Page:       "pages/index/index",
		CheckPath:  true,
		EnvVersion: "release",
		Width:      430,
		AutoColor:  false,
		LineColor: Color{
			G: 0,
			R: 0,
			B: 0,
		},
		IsHyaline: false,
	}
	for _, set := range setters {
		set(&options)
	}
	token, err := q.auth.GetAccessToken()
	if err != nil {
		return
	}
	var s string
	for key, value := range scene {
		s += fmt.Sprintf("&%s=%s", key, value)
	}
	options.Scene = strings.TrimLeft(s, "&")
	resp, err := http.PostJSON(fmt.Sprintf(qrCodeUnlimited, token), options)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &r)
	if err == nil {
		return
	} else {
		r.Content = resp
		err = nil
	}
	return
}
