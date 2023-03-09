package sns

import (
	"encoding/json"
	"fmt"

	"github.com/zedisdog/ty/sdk/net/http"
)

const (
	userinfoTemp = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s"
)

func NewSns() *Sns {
	return &Sns{
		Lang: "zh_CN",
	}
}

type Sns struct {
	Lang string
}

func (s Sns) UserInfo(userAccessToken string, openID string) (info UserInfoRes, err error) {
	url := fmt.Sprintf(userinfoTemp, userAccessToken, openID, s.Lang)
	resp, err := http.GetJSON(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &info)
	return
}
