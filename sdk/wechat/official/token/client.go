package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/sdk/net/http"
)

const (
	TokenTemp = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func NewToken(appID, secret string) *Token {
	return &Token{
		appID:  appID,
		secret: secret,
	}
}

type Token struct {
	appID           string
	secret          string
	lastAccessToken string
	lastTime        time.Time
	expiresIn       uint64
}

func (t *Token) GetAccessToken() (accessToken string, err error) {
	if t.lastAccessToken == "" || time.Since(t.lastTime) > (time.Duration(t.expiresIn)*time.Second) {
		var token AccessToken
		token, err = t.doGetAccessToken()
		if err != nil {
			err = errx.Wrap(err, "get access token from remote failed")
			return
		}
		t.lastAccessToken = token.AccessToken
		t.lastTime = time.Now()
		t.expiresIn = (token.ExpiresIn - 5)
	}
	accessToken = t.lastAccessToken
	return
}

func (t Token) doGetAccessToken() (accessToken AccessToken, err error) {
	url := fmt.Sprintf(TokenTemp, t.appID, t.secret)
	resp, err := http.GetJSON(url)
	if err != nil {
		err = errx.Wrap(err, "http error")
		return
	}
	err = json.Unmarshal(resp, &accessToken)
	if err != nil {
		err = errx.Wrap(err, "json invalid")
		return
	}

	if accessToken.ErrCode != 0 {
		err = errx.New(fmt.Sprintf("errcode: %d, errmsg: %s", accessToken.ErrCode, accessToken.ErrMsg))
		return
	}

	return
}
