package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义状态码
const (
	ERROR   = 0
	SUCCESS = 1
)

// Response 返回数据包装
type Rsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result 手动组装返回结果
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Rsp{
		code,
		msg,
		data,
	})
}

// FailWithMessage 返回自定义消息的失败
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}
