package controllers

import (
	"congchat-user/db"
	"congchat-user/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetMomentResponse struct {
	Moments []model.Moment `json:"moments"`
}

// 发送发布朋友圈动态请求
type CreateMomentRequest struct {
	Content string `json:"type:text" binding:"required"`
	UserID  uint   `gorm:"index"`
	ImgURL  string `gorm:"type:text"` //图片或视频地址。用于前端加载视频
	Goods   int    `gorm:"not null"`
	GoodsID []int  `gorm:"not null"`
}

// 查看朋友圈动态请求（）本质应该是一个列表（含有朋友动态）//滑动查询
type GetMomentRequest struct {
	StartIdx int //0  -  9
	Limit    int //10
	//条件
	FriendID  uint
	SearchStr string
}

// 10.21.01处理发布动态时刻的接口
func CreateMomentHandler(c *gin.Context) {
	var req CreateMomentRequest
	//绑定参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt64("user_id")
	//TODO 参数校验（是否符合需求，万一涉及暴力，黄色，广告。直接拦截）

	// 创建新的Moment并保存到数据库（校验）
	moment := &model.Moment{
		UserID:  uint(userID),
		Content: req.Content,
		ImgURL:  req.ImgURL,
		Goods:   req.Goods,
		GoodsID: req.GoodsID,
	}
	result := db.Db.Create(moment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create moment"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Moment created successfully", "moment_id": moment.ID})

}
func EditMomentHandler(c *gin.Context) {
	// 从URL参数中获取动态ID
	momentIDStr := c.Param("moment_id")
	momentID, err := strconv.Atoi(momentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid moment ID"})
	}

	//初始化更新的朋友圈（有待把数据放入originalmoment里）
	var updateMoment model.Moment
	if err := c.ShouldBind(&updateMoment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//初始化一个新的结构体来接受
	var originalMoment model.Moment

	result := db.Db.First(&originalMoment, momentID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Moment not found"})
		return
	}
	originalMoment.Content = updateMoment.Content
	originalMoment.ImgURL = updateMoment.ImgURL

	// 保存更改到数据库
	if err := db.Db.Save(&originalMoment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update moment"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Moment updated successfully"})

}

// 10.27.01DeleteMomentHandler处理删除朋友圈动态函数
func DeleteMomentHandler(c *gin.Context) {
	// 从URL参数中获取动态ID
	momentIDStr := c.Param("moment_id")
	momentID, err := strconv.Atoi(momentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid moment ID"})
		return
	}

	//初始化moment动态 根据momentid来删除动态
	var moment model.Moment
	//绑参
	result := db.Db.Where("id = ?", momentID).Delete(&moment)
	//校验
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete moment"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Moment not found"})
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Moment deleted successfully"})

}

// 10.21.01处理查看朋友圈动态时刻的接口
func GetMomentHandler(c *gin.Context) {
	var req GetMomentRequest //初始化朋友圈动态请求

	//绑参
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 构建查询（校验）
	var moments []model.Moment

	// 限制Limit的最大值为10
	if req.Limit > 10 {
		req.Limit = 10
	}

	// 预处理分页参数
	offset := req.StartIdx

	// 初始化查询条件
	var whereConditions []interface{}
	var whereArgs []interface{}
	if req.FriendID != 0 {
		whereConditions = append(whereConditions, "user_id = ?")
		whereArgs = append(whereArgs, req.FriendID)
	}
	if req.SearchStr != "" {
		searchStr := "%" + req.SearchStr + "%"
		whereConditions = append(whereConditions, "content LIKE ?")
		whereArgs = append(whereArgs, searchStr)
	}

	query := db.Db.Model(&model.Moment{}).Preload("Comments").Order("created_at DESC")
	if len(whereConditions) > 0 {
		query = query.Where(whereConditions, whereArgs...)
	}
	err := query.Offset(offset).Limit(req.Limit).Find(&moments).Error

	// 处理查询结果
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch moments"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "GetMoment successfully", "data": moments})
}
