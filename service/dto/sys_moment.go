package dto

import "congchat-user/model"

// CreateMomentRequest 发送发布朋友圈动态请求
type CreateMomentRequest struct {
	Content string `json:"type:text" binding:"required"`
	UserID  uint   `gorm:"index"`
	ImgURL  string `gorm:"type:text"` //图片或视频地址。用于前端加载视频
	Goods   int    `gorm:"not null"`
	GoodsID []int  `gorm:"not null"`
}

// EditMomentRequest 编辑朋友圈动态请求
type EditMomentRequest struct {
	model.Moment
}

// GetMomentRequest 查看朋友圈动态请求（）本质应该是一个列表（含有朋友动态）//滑动查询
type GetMomentRequest struct {
	StartIdx int `json:"start_idx,omitempty"` //0  -  9
	Limit    int `json:"limit,omitempty"`     //10
	//条件
	FriendID  uint   `json:"friend_id,omitempty"`
	SearchStr string `json:"search_str,omitempty"`
}

// DeleteMomentRequest 删除朋友圈动态请求
type DeleteMomentRequest struct {
	Ids []int `json:"ids"`
}
