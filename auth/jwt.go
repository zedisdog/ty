package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenIsInvalid = errors.New("token is invalid")
)

const (
	JwtID        = "jti"
	JwtIssuer    = "iss"
	JwtSubject   = "sub"
	JwtAudience  = "aud"
	JwtExpiresAt = "exp"
	JwtNotBefore = "nbf"
	JwtIssuedAt  = "iat"
)

func NewJwtTokenBuilder() *jwtTokenBuilder {
	return &jwtTokenBuilder{
		method: jwt.SigningMethodHS256,
	}
}

type jwtTokenBuilder struct {
	jwt.MapClaims
	key    []byte
	method jwt.SigningMethod
}

func (j *jwtTokenBuilder) WithClaim(key string, value interface{}) *jwtTokenBuilder {
	if j.MapClaims == nil {
		j.MapClaims = make(jwt.MapClaims)
	}
	j.MapClaims[key] = value
	return j
}

func (j *jwtTokenBuilder) WithClaims(claims map[string]interface{}) *jwtTokenBuilder {
	if j.MapClaims == nil {
		j.MapClaims = make(jwt.MapClaims)
	}
	for key, value := range claims {
		j.MapClaims[key] = value
	}
	return j
}

func (j *jwtTokenBuilder) WithKey(key string) *jwtTokenBuilder {
	j.key = []byte(key)
	return j
}

func (j *jwtTokenBuilder) WithMethod(method jwt.SigningMethod) *jwtTokenBuilder {
	j.method = method
	return j
}

func (j *jwtTokenBuilder) BuildToken() (string, error) {
	t := jwt.NewWithClaims(j.method, j.MapClaims)
	return t.SignedString([]byte(j.key))
}

// Parse token and fill claims
func Parse(token string, key []byte) (t *jwt.Token, err error) {
	t, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	return
}
