package dto

import "errors"

type AddFriendRequest struct {
	FriendUsername string `json:"friend_username" binding:"required"`
	UserID         uint   `json:"user_id" binding:"required,gte=1"` // 用户ID必填
}

type SearchFriendRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`   // 用户ID必填
	UserName string `json:"user_name" binding:"required"` //用户名称
	NickName string `json:"nick_name" binding:"required"` // 用户的昵称
}

type AcceptFriendRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 申请用户ID
}

type RejectFriendRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 申请用户ID
}

// AddFriendClear 检查 FriendUsername 和 UserID 和 MomentID是否都不为零或空
func AddFriendClear(d *AddFriendRequest) error {
	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	//检查FriendUsername是否为空
	if d.FriendUsername == "" {
		return errors.New("好友用户名不能为空")
	}

	return nil
}

// SearchFriendClear 检查 UserName 和 UserID 是否都不为零或空
func SearchFriendClear(d *SearchFriendRequest) error {
	//检查UserName是否为空
	if d.UserName == "" {
		return errors.New("用户名不能为空")
	}

	//检查UserID是否为空
	if d.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	return nil
}
