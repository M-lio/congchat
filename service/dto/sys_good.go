package dto

type AddGoodRequest struct {
	MomentID uint `json:"moment_id" binding:"required"`
	UserID   uint `json:"user_id" binding:"required"` // 用户ID必填
}

type CancelGoodRequest struct {
	UserID   uint `json:"user_id" binding:"required"`   // 用户ID必填
	MomentID uint `json:"moment_id" binding:"required"` //用户名称
}
