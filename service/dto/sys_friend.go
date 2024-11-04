package dto

type AddFriendRequest struct {
	FriendUsername string `json:"friend_username" binding:"required"`
	UserID         uint   `json:"user_id" binding:"required"` // 用户ID必填
}

type SearchFriendRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`   // 用户ID必填
	UserName string `json:"user_name" binding:"required"` //用户名称
	NickName string `json:"nick_name" binding:"required"` // 用户的昵称
}
