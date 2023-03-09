package ticket

import "github.com/zedisdog/ty/sdk/wechat/official/response"

type TicketResp struct {
	response.Error
	Ticket    string `json:"ticket"`
	ExpiresIn uint64 `json:"expires_in"`
}
