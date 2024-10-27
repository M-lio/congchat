package route

import (
	"congchat-user/controllers"
	"github.com/gin-gonic/gin"
)

// 修改时间10.24

// initUserRoute 初始化登录路由信息
func initUserRoute(r *gin.Engine) {
	Group := r.Group("/user")

	// 获取用户资料(其实也要先进行身份验证吗?)//修改时间10.24.1
	Group.GET("/:id", controllers.GetUserHandle)

	// 更新用户资料()//修改时间10.24.2
	Group.PUT("/:id", controllers.UpdateUserHandle)

	//发送获取好友列表请求
	Group.GET("/friend", controllers.GetFriendsHandler)
}
