package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// initRegisterRoute 初始化登录路由信息
func initRegisterRoute(r *gin.Engine) {
	// 定义公开的路由（不受身份验证保护）
	r.GET("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a register endpoint."})
	})

}

// initLoginRoute 初始化登录路由信息
func initLoginRoute(r *gin.Engine) {
	// 获取登录
	r.GET("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a login endpoint."})
	})
}
