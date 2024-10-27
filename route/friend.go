package route

import (
	"congchat-user/controllers"
	"congchat-user/middleware"
	"github.com/gin-gonic/gin"
)

// initFriendRoute 初始化登录路由信息
func initFriendRoute(r *gin.Engine) {
	// 添加搜索好友的路由1
	r.GET("/search/friends", controllers.SearchFriendsHandler)

	//发送好友请求2
	r.POST("/add-friend", controllers.AddFriendHandler)

	//接受好友请求3
	r.POST("/friendships/:id/accept", middleware.AuthMiddleware(controllers.AcceptFriendsRequestHandler))

	//拒绝好友请求4
	r.POST("/friendships/:id/reject", middleware.AuthMiddleware(controllers.RejectFriendsRequestHandler))
}
