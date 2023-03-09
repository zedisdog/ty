package official

import (
	"github.com/zedisdog/ty/sdk/wechat/official/jsconfig"
	"github.com/zedisdog/ty/sdk/wechat/official/oauth"
	"github.com/zedisdog/ty/sdk/wechat/official/sns"
	"github.com/zedisdog/ty/sdk/wechat/official/ticket"
	"github.com/zedisdog/ty/sdk/wechat/official/token"
)

func NewClient(appID string, secret string) *Client {
	return &Client{
		appID:  appID,
		secret: secret,
	}
}

type Client struct {
	appID    string
	secret   string
	oauth    *oauth.Oauth
	sns      *sns.Sns
	token    *token.Token
	ticket   *ticket.Ticket
	jsConfig *jsconfig.JsConfig
}

func (c *Client) Oauth() *oauth.Oauth {
	if c.oauth == nil {
		c.oauth = oauth.NewOauth(c.appID, c.secret)
	}
	return c.oauth
}

func (c *Client) Sns() *sns.Sns {
	if c.sns == nil {
		c.sns = sns.NewSns()
	}
	return c.sns
}

func (c *Client) Token() *token.Token {
	if c.token == nil {
		c.token = token.NewToken(c.appID, c.secret)
	}
	return c.token
}

func (c *Client) Ticket() *ticket.Ticket {
	if c.ticket == nil {
		c.ticket = ticket.NewTiket(c.Token())
	}
	return c.ticket
}

func (c *Client) JsConfig() *jsconfig.JsConfig {
	if c.jsConfig == nil {
		c.jsConfig = jsconfig.NewJsConfig(c.appID, c.Ticket())
	}

	return c.jsConfig
}
