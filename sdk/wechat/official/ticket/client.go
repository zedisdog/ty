package ticket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/wechat/official/token"
)

const (
	getTieketTemp = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

func NewTiket(token *token.Token) *Ticket {
	return &Ticket{
		token: token,
	}
}

type Ticket struct {
	lastTicket string
	lastTime   time.Time
	expiresIn  uint64
	token      *token.Token
}

func (t *Ticket) GetTicket() (ticket string, err error) {
	if t.lastTicket == "" || time.Since(t.lastTime) > (time.Duration(t.expiresIn)*time.Second) {
		var ticketResp TicketResp
		ticketResp, err = t.doGetJsTicket()
		if err != nil {
			return
		}
		t.lastTicket = ticketResp.Ticket
		t.lastTime = time.Now()
		t.expiresIn = (ticketResp.ExpiresIn - 5)
	}
	ticket = t.lastTicket
	return
}

func (t Ticket) doGetJsTicket() (ticket TicketResp, err error) {
	accessToken, err := t.token.GetAccessToken()
	if err != nil {
		err = errx.Wrap(err, "get ticket remote failed")
		return
	}
	url := fmt.Sprintf(getTieketTemp, accessToken)
	resp, err := http.GetJSON(url)
	if err != nil {
		err = errx.Wrap(err, "http error")
		return
	}
	err = json.Unmarshal(resp, &ticket)
	if err != nil {
		err = errx.Wrap(err, "invalid json")
	}

	if ticket.ErrCode != 0 {
		err = errx.New(fmt.Sprintf("errcode: %d, errmsg: %s", ticket.ErrCode, ticket.ErrMsg))
		return
	}

	return
}
