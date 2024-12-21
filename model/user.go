package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"` // 用户密码，这里不存储明文密码
	//OpenID    string    `json:"open_id" gorm:"uniqueIndex"` // 用户的OpenID，唯一标识
	Avatar   string // 用户头像URL
	NickName string `json:"nick_name"` // 用户的昵称
	Gender   int    `json:"gender"`    // 用户性别，0：未知、1：男性、2：女性
	Balance  float64
}
