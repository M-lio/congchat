package db

import (
	"congchat-user/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

// InitDB 初始化连接数据库

func InitDB() {
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// 从环境变量中读取数据库密码
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatalf("DB_PASSWORD environment variable is not set")
	}

	// 构造数据源名称（DSN）
	dsn := fmt.Sprintf("root:%s@tcp(z1.juhong.live:3306)/congchat_user?charset=utf8mb4&parseTime=True&loc=Local", dbPassword)

	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移模式，确保 用户User存在
	err = Db.AutoMigrate(&model.User{})
	err = Db.AutoMigrate(&model.Transfer{})
	err = Db.AutoMigrate(&model.Friendship{})
	//err = Db.AutoMigrate(&model.Moment{})
	err = Db.AutoMigrate(&model.Goods{})
	//err = Db.AutoMigrate(&model.Comment{})
	if err != nil {
		panic(err)
	}
}
