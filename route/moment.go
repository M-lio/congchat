package route

import (
	"congchat-user/controllers"
	"congchat-user/middleware"
	"github.com/gin-gonic/gin"
)

// initMomentsRoute 初始化登录路由信息
func initMomentsRoute(r *gin.Engine) {

	// 发布朋友圈Moments的请求路由5
	r.POST("/moments", middleware.AuthMiddleware(controllers.CreateMomentHandler))

	//删除朋友圈Moments的请求路由6
	r.DELETE("/moments/:moment_id", controllers.DeleteMomentHandler)

	//编辑朋友圈Moments的请求路由7
	r.PUT("/moments/:moment_id", controllers.EditMomentHandler)

	//查看朋友圈的路由8
	r.GET("/get-moments", middleware.AuthMiddleware(controllers.GetMomentHandler))

	//点赞朋友圈9
	r.POST("/goods", controllers.GoodsMomentHandler)

	//取消点赞朋友圈10
	r.DELETE("/goods", controllers.CancelGoodsMomentHandler)

	// 添加评论路由 11
	r.POST("/moments/:moment_id/comments", controllers.CommentHandler)

	// 删除评论路由  12
	r.DELETE("/comments/:id", controllers.DeleteCommentHandler)
}
