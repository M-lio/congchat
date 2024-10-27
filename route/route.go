package route

import "github.com/gin-gonic/gin"

// InitRoute 初始化路由
func InitRoute(r *gin.Engine) {

	//初始化注册路由
	initRegisterRoute(r)

	// 初始化登录相关路由
	initLoginRoute(r)

	// 初始化用户相关路由
	initUserRoute(r)

	// 初始化好友相关路由
	initFriendRoute(r)

	// 初始化朋友圈相关路由
	initMomentsRoute(r)

	// 初始化消息模块路由
	//initMessageRoute(app)
}

//只有route.route.go这页里面的逻辑处理函数是大写代表最高公民InitRoute 管理着所有的路由初始化函数initRoute
