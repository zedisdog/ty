package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/wechat/mini/auth/cache"
	"github.com/zedisdog/ty/sdk/wechat/mini/auth/request"
	"github.com/zedisdog/ty/sdk/wechat/mini/auth/response"
)

const (
	getAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	getPhoneNumberUrl = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
)

func WithCache(cache cache.Cache) func(*Auth) {
	return func(a *Auth) {
		a.cache = cache
	}
}

func NewAuth(appID string, secret string, setters ...func(*Auth)) *Auth {
	a := &Auth{
		appID:  appID,
		secret: secret,
		cache:  new(cache.DefaultCache),
	}
	for _, set := range setters {
		set(a)
	}
	return a
}

type Auth struct {
	appID  string
	secret string
	cache  cache.Cache
}

func (a Auth) GetAccessToken() (token string, err error) {
	token = a.cache.GetAccessToken()
	if token != "" {
		return token, nil
	}

	r, err := a.DoGetAccessToken()
	if err != nil {
		return
	}
	a.cache.SetAccessToken(r.AccessToken, r.ExpiresIn)
	return r.AccessToken, nil
}

func (a Auth) DoGetAccessToken() (r response.AccessToken, err error) {
	resp, err := http.GetJSON(fmt.Sprintf(getAccessTokenUrl, a.appID, a.secret))
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &r)
	return
}

func (a Auth) GetPhoneNumber(code string) (r response.PhoneNumber, err error) {
	token, err := a.GetAccessToken()
	if err != nil {
		return
	}
	resp, err := http.PostJSON(fmt.Sprintf(getPhoneNumberUrl, token), request.PhoneNumber{
		Code: code,
	})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &r)
	return
}
