package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/wechat/official/response"
)

const (
	SnsapiBase       = "snsapi_base"
	SnsapiUserinfo   = "snsapi_userinfo"
	oauthUrlTmpl     = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	code2AccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

func NewOauth(appID string, secret string) *Oauth {
	return &Oauth{
		appID:  appID,
		secret: secret,
	}
}

type Oauth struct {
	appID  string
	secret string
}

func (o Oauth) GenRedirectUrl(callbackUrl string, options ...func(*redirectOptions)) string {
	opt := &redirectOptions{
		scope: SnsapiBase,
	}
	for _, set := range options {
		set(opt)
	}
	return fmt.Sprintf(oauthUrlTmpl, o.appID, callbackUrl, opt.scope, opt.State())
}

func (o Oauth) Code2AccessToken(code string) (res response.Auth2AccessTokenRes, err error) {
	url := fmt.Sprintf(code2AccessToken, o.appID, o.secret, code)
	resp, err := http.GetJSON(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &res)
	return
}
