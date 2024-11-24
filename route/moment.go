package route

import (
	"congchat-user/controllers"
	"congchat-user/middleware"
	"github.com/gin-gonic/gin"
)

// initMomentsRoute 初始化登录路由信息
func initMomentsRoute(r *gin.Engine) {

	MomentApi := controllers.SysMoment{}

	GoodApi := controllers.SysGoods{}

	CommentApi := controllers.SysComment{}

	// 发布朋友圈Moments的请求路由5
	r.POST("/moments", middleware.AuthMiddleware(MomentApi.Insert))

	//删除朋友圈Moments的请求路由6
	r.DELETE("/moments/:moment_id", MomentApi.Delete)

	//编辑朋友圈Moments的请求路由7
	r.PUT("/moments/:moment_id", MomentApi.Edit)

	//查看朋友圈的路由8
	r.GET("/get-moments", middleware.AuthMiddleware(MomentApi.Get))

	//点赞朋友圈9
	r.POST("/goods", GoodApi.Add)

	//取消点赞朋友圈10
	r.DELETE("/goods", GoodApi.Cancel)

	// 添加评论路由 11
	r.POST("/moments/:moment_id/comments", CommentApi.Insert)

	// 删除评论路由  12
	r.DELETE("/comments/:id", CommentApi.Delete)
}
