package core

import "github.com/gin-gonic/gin"

// NotFoundErrorHandler 404异常处理
func NotFoundErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		FailWithMessage("错误路径", ctx)
		ctx.Abort()
	}
}
