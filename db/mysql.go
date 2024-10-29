package db

import (
	"congchat-user/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 初始化数据库连接
func InitDB() {
	dsn := "root:52Tiananmen.@tcp(z1.juhong.live:3306)/congchat_user?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移模式，确保 用户User存在
	err = Db.AutoMigrate(&model.User{})
	err = Db.AutoMigrate(&model.Moment{})
	err = Db.AutoMigrate(&model.Goods{})
	err = Db.AutoMigrate(&model.Comment{})
	if err != nil {
		panic(err)
	}

}
