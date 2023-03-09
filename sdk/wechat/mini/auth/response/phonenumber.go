package response

import "github.com/zedisdog/ty/sdk/wechat/mini/common"

type WaterMark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

type PhoneInfo struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       WaterMark `json:"watermark"`
}

type PhoneNumber struct {
	common.ErrorResponse
	PhoneInfo PhoneInfo `json:"phone_info"`
}
