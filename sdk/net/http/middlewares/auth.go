package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zedisdog/ty/sdk/net/http/response"
	"github.com/zedisdog/ty/slice"

	"github.com/zedisdog/ty/auth"

	"github.com/gin-gonic/gin"
)

func NewAuthBuilder() *authBuilder {
	return &authBuilder{
		userIdentityFrom: []string{auth.JwtSubject},
		tokenIDFrom:      auth.JwtID,
		cacheClaims:      true,
		tokenHeaderField: "Authorization",
		strict:           true,
		onTokenInvalid: func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "bearer token is invalid"})
		},
		onNoTokenFound: func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "no token found"})
		},
		onTokenParseFailed: func(ctx *gin.Context, err error) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "token is invalid1"})
		},
		onClaimsInvalid: func(ctx *gin.Context, claims interface{}) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid2",
			})
		},
		onIDNotFound: func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid3",
			})
		},
		onIDParseFailed: func(ctx *gin.Context, id interface{}) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid4",
			})
		},
		onFindUserFailed: func(ctx *gin.Context, err error) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid5",
			})
		},
		onUserNotExists: func(ctx *gin.Context, id uint64) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid6",
			})
		},
		onPassFailed: func(ctx *gin.Context, err error) {
			response.Error(ctx, err, http.StatusUnauthorized)
		},
	}
}

// authBuilder auth middleware builder. it parses token with given conditions.
type authBuilder struct {
	userIdentityFrom   []string                      //field name of user identity in token
	tokenIDFrom        string                        //filed name of token identity in token
	userExists         func(id uint64) (bool, error) //function to determine if user exists
	authKey            []byte                        //salt used by generate jwt signature
	cacheClaims        bool                          //if cache claims into context
	onPass             func(claims jwt.MapClaims, ctx *gin.Context) error
	tokenHeaderField   string
	strict             bool
	onTokenInvalid     func(ctx *gin.Context)
	onNoTokenFound     func(ctx *gin.Context)
	onTokenParseFailed func(ctx *gin.Context, err error)
	onClaimsInvalid    func(ctx *gin.Context, claims interface{})
	onIDNotFound       func(ctx *gin.Context)
	onIDParseFailed    func(ctx *gin.Context, id interface{})
	onFindUserFailed   func(ctx *gin.Context, err error)
	onUserNotExists    func(ctx *gin.Context, id uint64)
	onPassFailed       func(ctx *gin.Context, err error)
	whiteList          []uint64
}

func (ab *authBuilder) WithUserIdentityFrom(jwtField ...string) *authBuilder {
	ab.userIdentityFrom = jwtField
	return ab
}

func (ab *authBuilder) WithTokenIDFrom(jwtField string) *authBuilder {
	ab.tokenIDFrom = jwtField
	return ab
}

func (ab *authBuilder) WithOnPass(f func(claims jwt.MapClaims, ctx *gin.Context) error) *authBuilder {
	ab.onPass = f
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

func (ab *authBuilder) WithTokenHeaderField(field string) *authBuilder {
	ab.tokenHeaderField = field
	return ab
}

func (ab *authBuilder) WithStrict(enable bool) *authBuilder {
	ab.strict = enable
	return ab
}

func (ab *authBuilder) WithOnTokenInvalid(f func(ctx *gin.Context)) *authBuilder {
	ab.onTokenInvalid = f
	return ab
}

func (ab *authBuilder) WithOnNoTokenFound(f func(ctx *gin.Context)) *authBuilder {
	ab.onNoTokenFound = f
	return ab
}

func (ab *authBuilder) WithOnTokenParseFailed(f func(ctx *gin.Context, err error)) *authBuilder {
	ab.onTokenParseFailed = f
	return ab
}

func (ab *authBuilder) WithOnClaimsInvalid(f func(ctx *gin.Context, claims interface{})) *authBuilder {
	ab.onClaimsInvalid = f
	return ab
}

func (ab *authBuilder) WithOnIDNotFound(f func(ctx *gin.Context)) *authBuilder {
	ab.onIDNotFound = f
	return ab
}

func (ab *authBuilder) WithOnIDParseFailed(f func(ctx *gin.Context, id interface{})) *authBuilder {
	ab.onIDParseFailed = f
	return ab
}

func (ab *authBuilder) WithOnFindUserFailed(f func(ctx *gin.Context, err error)) *authBuilder {
	ab.onFindUserFailed = f
	return ab
}

func (ab *authBuilder) WithOnUserNotExists(f func(ctx *gin.Context, id uint64)) *authBuilder {
	ab.onUserNotExists = f
	return ab
}

func (ab *authBuilder) WithOnPassFailed(f func(ctx *gin.Context, err error)) *authBuilder {
	ab.onPassFailed = f
	return ab
}

func (ab *authBuilder) WithPassIDs(ids []uint64) *authBuilder {
	ab.whiteList = ids
	return ab
}

func (ab *authBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		if ctx.Request.Header.Get(ab.tokenHeaderField) != "" {
			arr := strings.Split(ctx.Request.Header.Get(ab.tokenHeaderField), " ")
			switch len(arr) {
			case 0:
				ab.onTokenInvalid(ctx)
				return
			case 1:
				if ab.strict {
					ab.onTokenInvalid(ctx)
					return
				}
				token = arr[0]
			default:
				token = arr[1]
			}
		} else if token, _ = ctx.Cookie("token"); token != "" {
		} else if ctx.Query("token") != "" {
			token = ctx.Query("token")
		} else {
			println("no token found")
			ab.onNoTokenFound(ctx)
			return
		}

		t, err := auth.Parse(token, ab.authKey)
		if err != nil || !t.Valid {
			println("parse failed")
			ab.onTokenParseFailed(ctx, err)
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			println("parse failed2")
			ab.onClaimsInvalid(ctx, t.Claims)
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
				ab.onIDNotFound(ctx)
				return
			}

			var id uint64
			id, err = strconv.ParseUint(IDStr.(string), 10, 64)
			if err != nil {
				ab.onIDParseFailed(ctx, IDStr)
				return
			}

			if ab.userExists != nil {
				exists, err := ab.userExists(id)
				if err != nil {
					ab.onFindUserFailed(ctx, err)
					return
				}
				if !exists && !slice.Containers(id, ab.whiteList) {
					ab.onUserNotExists(ctx, id)
					return
				}
			}

			ctx.Set("id", id)
		}

		if ab.onPass != nil {
			err = ab.onPass(claims, ctx)
			if err != nil {
				ab.onPassFailed(ctx, err)
				return
			}
		}

		if ab.cacheClaims {
			ctx.Set("claims", claims)
		}

		ctx.Next()
	}
}
