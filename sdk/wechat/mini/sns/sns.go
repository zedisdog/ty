package sns

import (
	"encoding/json"
	"fmt"

	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/wechat/mini/sns/response"
)

const (
	code2sessionTmpl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func NewSns(appID string, secret string) *Sns {
	return &Sns{
		appID:  appID,
		secret: secret,
	}
}

type Sns struct {
	appID  string
	secret string
}

func (s Sns) Code2Session(code string) (r response.Code2SessionResponse, err error) {
	resp, err := http.GetJSON(fmt.Sprintf(code2sessionTmpl, s.appID, s.secret, code))
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &r)
	return
}
