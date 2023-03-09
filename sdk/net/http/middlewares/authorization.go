package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Checker func(uint64, ...interface{}) (bool, error)

func GenChecker(checker Checker, targets ...interface{}) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		pass, err := checker(ctx.MustGet("id").(uint64), targets...)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		if !pass {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{"message": "未授权的访问"})
			return
		}
		ctx.Next()
	}
}
