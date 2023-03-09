package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"

	"github.com/zedisdog/ty/auth"

	"github.com/gin-gonic/gin"
)

func NewAuthBuilder() *authBuilder {
	return &authBuilder{}
}

// authBuilder auth middleware builder. it parses token with given conditions.
type authBuilder struct {
	userIdentityFrom []string                      //field name of user identity in token
	tokenIDFrom      string                        //filed name of token identity in token
	roleFrom         string                        //filed name of role name in token
	userExists       func(id uint64) (bool, error) //function to determine if user exists
	authKey          []byte                        //salt used by generate jwt signature
	cacheClaims      bool                          //if cache claims into context
}

func (ab *authBuilder) WithUserIdentityFrom(jwtField ...string) *authBuilder {
	ab.userIdentityFrom = jwtField
	return ab
}

func (ab *authBuilder) WithTokenIDFrom(jwtField string) *authBuilder {
	ab.tokenIDFrom = jwtField
	return ab
}

func (ab *authBuilder) WithRoleFrom(jwtField string) *authBuilder {
	ab.roleFrom = jwtField
	return ab
}

func (ab *authBuilder) WithUserExistsFunc(f func(id uint64) (bool, error)) *authBuilder {
	ab.userExists = f
	return ab
}

func (ab *authBuilder) WithAuthKey(key string) *authBuilder {
	ab.authKey = []byte(key)
	return ab
}

func (ab *authBuilder) WithClaimsCache() *authBuilder {
	ab.cacheClaims = true
	return ab
}

func (ab *authBuilder) Build() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var token string
		if ctx.Request.Header.Get("Authorization") != "" {
			arr := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
			if len(arr) < 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "bearer token is invalid"})
				return
			}
			token = arr[1]
		} else if ctx.Query("token") != "" {
			token = ctx.Query("token")
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "no token found"})
			return
		}

		t, err := auth.Parse(token, ab.authKey)
		if err != nil || !t.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "token is invalid1"})
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid2",
			})
			return
		}

		if len(ab.userIdentityFrom) > 0 {
			var IDStr interface{}
			for _, field := range ab.userIdentityFrom {
				IDStr, ok = claims[field]
				if ok {
					break
				}
			}
			if IDStr == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid3",
				})
				return
			}

			var id uint64
			id, err = strconv.ParseUint(IDStr.(string), 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid4",
				})
				return
			}

			if ab.userExists != nil {
				exists, err := ab.userExists(id)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "token is invalid5",
					})
					return
				}
				if !exists {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "token is invalid6",
					})
					return
				}
			}

			ctx.Set("id", id)
		}

		if ab.roleFrom != "" {
			var role interface{}
			role, ok = claims[ab.roleFrom]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid7",
				})
				return
			}
			ctx.Set("role", role.(string))
		}

		if ab.cacheClaims {
			ctx.Set("claims", claims)
		}

		ctx.Next()
	}
}