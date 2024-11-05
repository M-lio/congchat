package controllers

import (
	"congchat-user/core"
	"congchat-user/model"
	"congchat-user/service"
	"congchat-user/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SysMoment struct {
	core.Api
}

type GetMomentResponse struct {
	Moments []model.Moment `json:"moments"`
}

// GetMomentRequest 查看朋友圈动态请求（）本质应该是一个列表（含有朋友动态）//滑动查询
type GetMomentRequest struct {
	StartIdx int //0  -  9
	Limit    int //10
	//条件
	FriendID  uint
	SearchStr string
}

// Insert 11.5处理发布动态时刻的接口
func (e SysMoment) Insert(c *gin.Context) {
	req := dto.CreateMomentRequest{}
	var rsp core.Rsp
	s := new(service.SysMoment)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.CreateMoment(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "朋友圈发布成功"
	c.JSON(http.StatusOK, rsp)
}

// 10.21.01处理发布动态时刻的接口旧代码
/*
func CreateMomentHandler(c *gin.Context) {
	var req dto.CreateMomentRequest
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

*/

// Edit 11.5处理编辑动态时刻的接口
func (e SysMoment) Edit(c *gin.Context) {
	req := dto.EditMomentRequest{}
	var rsp core.Rsp
	s := new(service.SysMoment)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.EditMoment(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "朋友圈编辑成功"
	c.JSON(http.StatusOK, rsp)
}

/*
func EditMomentHandler(c *gin.Context) {
	// 从URL参数中获取动态ID
	momentIDStr := c.Param("moment_id")
	momentID, err := strconv.Atoi(momentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid moment ID"})
	}

	//初始化更新的朋友圈（有待把数据放入originalMoment里）
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

*/ // 处理编辑动态时刻的接口旧代码

// Delete 处理删除朋友圈动态函数
func (e SysMoment) Delete(c *gin.Context) {
	req := dto.DeleteMomentRequest{}
	var rsp core.Rsp
	s := new(service.SysMoment)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.RemoveMoment(&req).Error
	if err != nil {
		return
	}
	rsp.Code = 0
	rsp.Msg = "删除朋友圈成功"
	c.JSON(http.StatusOK, rsp)
}

// 10.27.01DeleteMomentHandler处理删除朋友圈动态函数旧代码
/*
func DeleteMomentHandler(c *gin.Context) {
	// 从URL参数中获取动态ID
	momentIDStr := c.Param("moment_id")
	momentID, err := strconv.Atoi(momentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid moment ID"})
		return
	}

	//初始化moment动态 根据momentID来删除动态
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
*/

// Get 处理查看朋友圈动态时刻的接口
func (e SysMoment) Get(c *gin.Context) {
	s := service.SysMoment{}
	var rsp core.Rsp
	req := dto.GetMomentRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var list = make([]model.Moment, 0)
	err := s.GetMoment(&req, &list).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp.Code = 0
	rsp.Data = list
	rsp.Msg = "获取朋友圈成功"
	c.JSON(http.StatusOK, rsp)
}

// 10.21.01处理查看朋友圈动态时刻的接口旧代码
/*
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
*/
