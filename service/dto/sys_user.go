package dto

type GetUserRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

type UpdateUserRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

type GetFriendsRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}
