package token

import "github.com/zedisdog/ty/sdk/wechat/official/response"

type AccessToken struct {
	response.Error
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
}
