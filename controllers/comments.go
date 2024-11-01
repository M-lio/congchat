package controllers

import (
	"congchat-user/core"
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service"
	"congchat-user/service/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
)

type SysComment struct {
	core.Api
}

//func (e SysMenu) GetPage(c *gin.Context) {
//	s := service.SysMenu{}  		生成服务实例（service）
//	req := dto.SysMenuGetPageReq{}	生成请求实例（请求在dto）
//	err := e.MakeContext(c).		设置上下文
//		MakeOrm().					生成DB
//		Bind(&req, binding.Form).	绑参校验
//		MakeService(&s.Service).	生成指定服务
//		Errors						错误记录
//	if err != nil {					看看执行中有没有错
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	var list = make([]models.SysMenu, 0) 执行中用列表记录错误信息
//	err = s.GetPage(&req, &list).Error
//	if err != nil {
//		e.Error(500, err, "查询失败")
//		return
//	}
//	e.OK(list, "查询成功")				返回
//}

type DeleteCommentResponse struct {
	Comments []model.Comment `json:"comments"`
}

type GetCommentsResponse struct {
	Comments []model.Comment `json:"comments"`
}

func (e SysComment) Insert(c *gin.Context) {
	req := dto.CreateCommentRequest{}
	var rsp core.Rsp
	s := new(service.SysComment)
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		// 使用 ShouldBindWith 而不是 BindJSON 可以让我们指定绑定类型，并且它会自动处理验证
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.CreateComment(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "评论成功"
	c.JSON(http.StatusOK, rsp)
}

func CreateCommentHandler(c *gin.Context) {
	var req dto.CreateCommentRequest
	var rsp core.Rsp
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

	rsp.Data = GetCommentsResponse{Comments: comments}
	rsp.Code = 0
	rsp.Msg = "评论成功"
	c.JSON(http.StatusOK, rsp)
}

func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("id")
	var rsp core.Rsp
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

	rsp.Code = 0
	rsp.Msg = "删除成功"
	c.JSON(http.StatusOK, rsp)
}
