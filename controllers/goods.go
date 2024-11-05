package controllers

import (
	"congchat-user/core"
	"congchat-user/service"
	"congchat-user/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SysGoods struct {
	core.Api
}

// Add 为指定用户给指定动态点赞  11.5
func (e SysGoods) Add(c *gin.Context) {
	var req dto.AddGoodRequest
	var rsp core.Rsp
	s := new(service.SysGoods)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.AddGood(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "点赞成功"
	c.JSON(http.StatusOK, rsp)
}

// GoodsMomentHandler 为指定用户给指定动态点赞旧代码
/*func GoodsMomentHandler(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id"`
		MomentID uint `json:"moment_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 检查是否已经点赞过
	var Good model.Goods
	result := db.Db.First(&Good, "user_id = ? AND moment_id = ?", request.UserID, request.MomentID)
	if result.Error != nil {
		// 处理数据库查询错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 { // RowsAffected 返回受影响的行数，如果没有找到则为0
		// 创建新的点赞记录
		newGood := model.Goods{UserID: request.UserID, MomentID: request.MomentID}
		if err := db.Db.Create(&newGood).Error; err != nil {
			// 处理创建记录错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
			return
		}

		// 更新朋友圈动态的点赞数
		var moment model.Moment
		if err := db.Db.First(&moment, request.MomentID).Error; err != nil {
			// 处理查询动态错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find moment"})
			return
		}
		moment.Goods++
		if err := db.Db.Save(&moment).Error; err != nil {
			// 处理更新点赞数错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update moment's goods count"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Good successfully"})
}

*/

// Cancel 取消指定用户对指定动态的点赞  11.5
func (e SysGoods) Cancel(c *gin.Context) {
	var req dto.CancelGoodRequest
	var rsp core.Rsp
	s := new(service.SysGoods)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.CancelGood(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "取消点赞成功"
	c.JSON(http.StatusOK, rsp)
}

// CancelGoodsMomentHandler 取消指定用户对指定动态的点赞 旧代码
/*
func CancelGoodsMomentHandler(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id"`
		MomentID uint `json:"moment_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 查找并删除点赞记录
	var Good model.Goods // 确保 Like 结构体中有 UserID 和 MomentID 字段
	result := db.Db.Delete(&Good, "user_id = ? AND moment_id = ?", request.UserID, request.MomentID)
	if result.Error != nil {
		// 处理数据库删除错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		// 如果没有找到要删除的记录，则返回错误或提示信息
		c.JSON(http.StatusNotFound, gin.H{"message": "Good not found"})
		return
	}

	// 更新朋友圈动态的点赞数
	var moment model.Moment
	if err := db.Db.First(&moment, request.MomentID).Error; err != nil {
		// 处理查询动态错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find moment"})
		return
	}
	// 假设 Goods 字段存储了点赞数，这里需要减一
	if moment.Goods > 0 {
		moment.Goods--
		if err := db.Db.Save(&moment).Error; err != nil {
			// 处理更新点赞数错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update moment's goods count"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Good canceled successfully"})
}

*/
