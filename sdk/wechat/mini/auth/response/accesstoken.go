package response

import "github.com/zedisdog/ty/sdk/wechat/mini/common"

type AccessToken struct {
	common.ErrorResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` //过期时间，单位秒
}
