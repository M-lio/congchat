package route

import (
	"congchat-user/controllers"
	"github.com/gin-gonic/gin"
)

// 修改时间10.24

// initUserRoute 初始化登录路由信息
func initUserRoute(r *gin.Engine) {
	UserApi := controllers.SysUser{}
	TransactionApi := controllers.SysTransferController{}

	Group := r.Group("/user")

	// 获取用户资料(其实也要先进行身份验证吗?)//修改时间10.24.1
	Group.GET("/:id", UserApi.Get)

	// 更新用户资料()//修改时间10.24.2
	Group.PUT("/:id", UserApi.Update)

	//发送获取好友列表请求
	Group.GET("/friend", UserApi.GetFriends)

	//// 发布朋友圈Moments的请求路由5
	r.POST("/moments", TransactionApi.Transfer)
}
