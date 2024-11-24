package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	// 其他必要的导入
)

// 地址：z1.juhong.live:6379
var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	// 设置Redis地址和密码（如果有）
	redisOptions := &redis.Options{
		Addr:     "z1.juhong.live:6379", // Redis服务器地址和端口
		Password: "52Tiananmen.",        // Redis密码（如果没有密码，可以为空）
		DB:       0,                     // 使用默认的Redis数据库索引
	}

	// 创建Redis客户端
	RedisClient = redis.NewClient(redisOptions)

	// 可选：测试连接是否成功
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
