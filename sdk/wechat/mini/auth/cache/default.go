package cache

import "time"

var _ Cache = (*DefaultCache)(nil)

type DefaultCache struct {
	token      string
	expireTime time.Time
}

func (d DefaultCache) GetAccessToken() string {
	if time.Now().Before(d.expireTime.Add(-5 * time.Second)) {
		return d.token
	} else {
		return ""
	}
}

func (d *DefaultCache) SetAccessToken(token string, expiresIn int) {
	d.token = token
	d.expireTime = time.Now().Add(time.Duration(expiresIn) * time.Second)
}
