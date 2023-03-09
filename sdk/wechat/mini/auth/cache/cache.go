package cache

type Cache interface {
	GetAccessToken() string
	SetAccessToken(token string, expiresIn int)
}
