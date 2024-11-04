package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"gorm.io/gorm"
)

type SysComment struct {
	Service
}

func (e *SysComment) CreateComment(c *dto.CreateCommentRequest) *SysComment {
	var err error
	// 检查是否已经存在相同的评论
	//todo 课后作业：了解mysql的事务。
	var existingComment model.Comment
	tx := e.Orm.Debug().Begin()
	err = db.Db.Where("moment_id = ? AND user_id = ? AND contents = ? ", c.MomentID, c.UserID, c.Contents).First(&existingComment).Error
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		_ = e.AddError(err)
		tx.Rollback()
	}
	err = db.Db.Create(&model.Comment{
		MomentID: c.MomentID,
		Contents: c.Contents,
		UserID:   c.UserID,
	}).Error
	if err != nil {
		_ = e.AddError(err)
	}

	//更新Moment的Comments数量
	var moment model.Moment
	err = db.Db.First(&moment, c.MomentID).Error
	if err != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	moment.Comments++

	db.Db.Save(&moment)
	tx.Commit()
	return e
}

// Remove 删除SysComment
func (e *SysComment) Remove(d *dto.DeleteCommentRequest) *SysComment {
	var err error
	var data model.Comment
	db.Db.First(&data, d.Ids)
	if err = db.Db.First(&data, d.Ids).Error; err != nil {
		_ = e.AddError(err)
	}

	if err = db.Db.Delete(&data, d.Ids).Error; err != nil {
		_ = e.AddError(err)
	}

	// 更新 Moment 的 Comments 数量（）
	var moment model.Moment
	if err = db.Db.First(&moment, data.MomentID).Error; err != nil {
		_ = e.AddError(err)
	}
	moment.Comments--
	if err = db.Db.Save(&moment).Error; err != nil {
		_ = e.AddError(err)
	}

	return e
}
