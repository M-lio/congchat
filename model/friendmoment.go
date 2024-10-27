package model

import (
	"gorm.io/gorm"
)

// Moment代表动态结构体
type Moment struct {
	gorm.Model
	UserID   uint   `gorm:"index"`
	Content  string `gorm:"type:text"`
	ImgURL   string `gorm:"type:text"` //图片或视频地址。用于前端加载视频
	Comment  []Comment
	Comments int   `gorm:"not null"`
	Goods    int   `gorm:"not null"`
	GoodsID  []int `gorm:"not null"`
}

// comment评论结构体
type Comment struct {
	gorm.Model
	MomentID   int    `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	Contents   string `gorm:"not null"`
	RelationID int    `gorm:"not null"` //是否回复。回复谁
}

//// Goods评论结构体
//type Goods struct {
//	gorm.Model
//	UserID
//}
