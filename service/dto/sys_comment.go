package dto

type CreateCommentRequest struct {
	Contents string `json:"contents" binding:"required,min=1,max=255"` // 假设内容至少1个字符，最多255个字符
	UserID   uint   `json:"user_id" binding:"required"`                // 用户ID必填
	MomentID uint   `json:"moment_id" binding:"required"`              // 动态ID必填
}

type DeleteCommentRequest struct {
	Ids []int `json:"ids"`
}
