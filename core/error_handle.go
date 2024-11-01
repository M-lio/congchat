package core

import (
	"github.com/gin-gonic/gin"
)

// NotFoundErrorHandler 404异常处理
func NotFoundErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		FailWithMessage("错误路径", ctx)
		ctx.Abort()
	}
}

// FailWithMessage 返回自定义消息的失败
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}
