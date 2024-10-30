package main

import (
	"congchat-user/core"
	"congchat-user/db"
	"congchat-user/route"
	"github.com/gin-gonic/gin"
	"log"
)

// TODO :权限常量用于IAM系统认证 有待开发
//const (
//	RoleUser  = "user"
//	RoleAdmin = "admin"
//)

func main() {
	db.InitDB()        // 初始化数据库连接
	r := gin.Default() // 创建一个默认的Gin引擎

	// 定义全局异常处理
	r.NoRoute(core.NotFoundErrorHandler())

	route.InitRoute(r) // 初始化登录相关的总路由

	//r.Use(LoggerMiddleware())//注册日志打印中间件

	// 启动服务器
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
