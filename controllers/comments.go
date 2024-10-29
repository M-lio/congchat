package controllers

import (
	"congchat-user/db"
	"congchat-user/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ComentHandler(c *gin.Context) {
	var comment model.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment payload"})
		return
	}

	result := db.Db.Create(&model.Comment{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Comment create failed"})
		return
	}

	//更新Moment的Comments数量
	var moment model.Moment
	if err := db.Db.First(&moment, comment.MomentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find moment"})
		return
	}
	moment.Comments++
	db.Db.Save(&moment)

	c.JSON(http.StatusOK, comment)
}

func DeleteComentHandler(c *gin.Context) {
	var comment model.Comment
	id := c.Param("id")
	if err := db.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
		return
	}

	result := db.Db.Delete(&model.Comment{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	//更新Moment的Comments数量
	moment := model.Moment{}
	if err := db.Db.First(&moment, comment.MomentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find moment"})
		return
	}

	moment.Comments--
	db.Db.Save(&moment)

	c.JSON(http.StatusOK, comment)
}
