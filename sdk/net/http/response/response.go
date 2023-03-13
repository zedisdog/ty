package response

import (
	"gorm.io/gorm"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zedisdog/ty/errx"
)

type Meta struct {
	CurrentPage uint `json:"current_page"`
	Total       uint `json:"total"`
	LastPage    uint `json:"last_page"`
	PerPage     uint `json:"per_page"`
}

type Response struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Meta *Meta       `json:"meta"`
}

// Error 返回错误响应 p1 错误 p2 status code
func Error(c *gin.Context, err error, status ...interface{}) {
	res := &Response{Msg: err.Error()}

	var code int
	if len(status) > 0 {
		code = status[0].(int)
	} else {
		if errx.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else if er, ok := err.(*errx.Error); ok && er.Code != 0 {
			code = int(er.Code)
			res.Data = er.Detail
		} else {
			code = http.StatusInternalServerError
		}
	}

	Json(c, res, code, false)
}

// Success params[0]: data
//
//	params[1]: status code
func Success(c *gin.Context, params ...interface{}) {
	var (
		code     int
		response *Response
	)

	switch len(params) {
	case 0:
		code = http.StatusNoContent
	case 1:
		response = &Response{
			Data: params[0],
		}
		code = http.StatusOK
	case 2:
		response = &Response{
			Data: params[0],
		}
		code = params[1].(int)
	}

	Json(c, response, code, false)
}

func Pagination(c *gin.Context, data interface{}, total int, page int, perPage int) {
	resp := &Response{
		Meta: &Meta{
			CurrentPage: uint(page),
			Total:       uint(total),
			LastPage:    uint(math.Ceil(float64(total) / float64(perPage))),
			PerPage:     uint(perPage),
		},
		Data: data,
	}
	if resp.Meta.CurrentPage == 0 {
		resp.Meta.CurrentPage = 1
	}
	if resp.Meta.LastPage == 0 {
		resp.Meta.LastPage = 1
	}
	Json(c, resp, http.StatusOK, false)
}

func Json(ctx *gin.Context, data interface{}, status int, abort bool) {
	if abort {
		ctx.AbortWithStatusJSON(status, data)
	} else {
		ctx.JSON(status, data)
	}
}
