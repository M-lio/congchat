package dto

import "errors"

type AddGoodRequest struct {
	MomentID uint `json:"moment_id" binding:"required"`
	UserID   uint `json:"user_id" binding:"required"` // 用户ID必填
}

type CancelGoodRequest struct {
	UserID   uint `json:"user_id" binding:"required"`   // 用户ID必填
	MomentID uint `json:"moment_id" binding:"required"` //用户名称
}

// AisClear 检查 MomentID 和 UserID 是否都不为零
func AisClear(d *AddGoodRequest) error {
	if d.MomentID == 0 {
		return errors.New("朋友圈ID不能为0")
	}
	if d.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	return nil
}

// CisClear 检查 MomentID 和 UserID 是否都不为零
func CisClear(d *CancelGoodRequest) error {
	if d.MomentID == 0 {
		return errors.New("朋友圈ID不能为0")
	}
	if d.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	return nil
}
