package response

import (
	"errors"
	"github.com/RyCarlos/go-common/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Response
// @Description 通用响应对象
type Response struct {
	// required:true
	// example: 0
	Code int `json:"code"` // 业务状态码
	// required:true
	// example: {}
	Data any `json:"data"` // 业务数据
	// required:true
	// example: 请求成功
	Msg string `json:"msg"` // 响应信息
	// required:true
	// example: 1751726010
	Timestamp int64 `json:"timestamp"` // 响应时间
}

type PageResponse[T any] struct {
	Total int64 `json:"total"`
	Items []*T  `json:"items"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		code,
		data,
		msg,
		time.Now().Unix(),
	})
}

func Ok(ctx *gin.Context) {
	Result(errs.NoneError, map[string]interface{}{}, "请求成功", ctx)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(errs.NoneError, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(errs.NoneError, data, "请求成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(errs.NoneError, data, message, c)
}

func Fail(ctx *gin.Context) {
	Result(errs.UnknownError, map[string]interface{}{}, "请求失败", ctx)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(errs.UnknownError, map[string]interface{}{}, message, c)
}

func FailWithError(err error, c *gin.Context) {
	var errCode errs.ErrCode
	code := errs.UnknownError
	if errors.As(err, &errCode) {
		code = errCode.Code()
	}
	Result(code, map[string]interface{}{}, err.Error(), c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(errs.UnknownError, data, message, c)
}

func NoRoute(c *gin.Context) {
	FailWithError(errs.ErrRoute, c)
}

func NoMethod(c *gin.Context) {
	FailWithError(errs.ErrHttpMethod, c)
}

func NoAuthWithError(err error, c *gin.Context) {
	FailWithError(err, c)
}
