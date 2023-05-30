package pkg

import (
	"github.com/xxandjg/ginbase/global"
	"net/http"
)

type Response struct {
	Code    int         `json:"code" xml:"code" yaml:"code" example:"10000"`
	Message string      `json:"message" xml:"message" yaml:"message" example:"操作成功"`
	Data    interface{} `json:"data" xml:"data" yaml:"data" example:"49ba59abbe56e057"`
}

type Context interface {
	JSON(code int, obj interface{})
	AbortWithStatusJSON(code int, jsonObj interface{})
}

func SuccessMsg(data interface{}) Response {
	result := Response{
		Code:    global.SUCCESS.GetCode(),
		Message: global.SUCCESS.GetMessage(),
		Data:    data,
	}
	return result
}
func ErrorMsg(err global.Error) Response {
	result := Response{
		Code:    err.GetCode(),
		Message: err.GetMessage(),
		Data:    nil,
	}
	return result
}

func Ok(ctx Context, data interface{}) {
	ctx.JSON(http.StatusOK, SuccessMsg(data))
}

func Err(ctx Context, err global.Error) {
	ctx.AbortWithStatusJSON(http.StatusOK, ErrorMsg(err))
}
