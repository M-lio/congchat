package controllers

import (
	"congchat-user/db"
	"congchat-user/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
)

type DeleteCommentResponse struct {
	Comments []model.Comment `json:"comments"`
}

type GetCommentsResponse struct {
	Comments []model.Comment `json:"comments"`
}

type CreateCommentRequest struct {
	Contents string `json:"contents" binding:"required,min=1,max=255"` // 假设内容至少1个字符，最多255个字符
	UserID   uint   `json:"user_id" binding:"required"`                // 用户ID必填
	MomentID uint   `json:"moment_id" binding:"required"`              // 动态ID必填
}

func CommentHandler(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		// 使用 ShouldBindWith 而不是 BindJSON 可以让我们指定绑定类型，并且它会自动处理验证
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}

	// 检查是否已经存在相同的评论
	var existingComment model.Comment
	r := db.Db.Where("moment_id = ? AND user_id = ? AND contents = ? ",
		req.MomentID, req.UserID, req.Contents).First(&existingComment)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate comment not allowed"})
			return
		}
	}

	result := db.Db.Create(&model.Comment{
		MomentID: req.MomentID,
		Contents: req.Contents,
		UserID:   req.UserID,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Comment create failed"})
		return
	}

	//更新Moment的Comments数量
	var moment model.Moment
	if err := db.Db.First(&moment, req.MomentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find moment"})
		return
	}
	moment.Comments++
	db.Db.Save(&moment)

	// 准备响应
	var comments []model.Comment
	db.Db.Preload("Moment").Where("moment_id = ?", req.MomentID).Find(&comments)
	if err := db.Db.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	getResponse := GetCommentsResponse{Comments: comments}
	c.JSON(http.StatusOK, getResponse)
}

func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("id")
	var comment model.Comment
	if err := db.Db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	result := db.Db.Delete(&comment, commentID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// 更新 Moment 的 Comments 数量（）
	var moment model.Moment
	if err := db.Db.First(&moment, comment.MomentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find moment"})
		return
	}
	moment.Comments--
	if err := db.Db.Save(&moment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Save moment"})
		return
	}

	// 获取并返回评论列表
	var comments []model.Comment
	db.Db.Preload("Moment").Where("moment_id = ?", comment.MomentID).Find(&comments)
	if err := db.Db.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	deleteResponse := DeleteCommentResponse{Comments: comments}
	c.JSON(http.StatusOK, deleteResponse)
}
