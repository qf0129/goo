package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type PageBody struct {
	List     any   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func Resp(c *gin.Context, httpStatus int, body RespBody) {
	c.JSON(httpStatus, body)
	c.Abort()
}

func RespOk(c *gin.Context, data any) {
	Resp(c, http.StatusOK, RespBody{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func RespFail(c *gin.Context, msg string) {
	Resp(c, http.StatusOK, RespBody{
		Code: 40000,
		Msg:  msg,
		Data: nil,
	})
}

func RespFailWith401(c *gin.Context, msg string) {
	Resp(c, http.StatusUnauthorized, RespBody{
		Code: 40100,
		Msg:  msg,
		Data: nil,
	})
}
