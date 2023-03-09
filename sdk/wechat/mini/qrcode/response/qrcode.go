package response

import "github.com/zedisdog/ty/sdk/wechat/mini/common"

type QrCodeUnlimited struct {
	common.ErrorResponse
	Content []byte
}
