package response

import "github.com/zedisdog/ty/sdk/wechat/mini/common"

type Code2SessionResponse struct {
	common.ErrorResponse
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}
