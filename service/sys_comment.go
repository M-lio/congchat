package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Msg   string
	MsgID string
	Error error
}

type SysComment struct {
	Service
}

func (e *SysComment) CreateComment(c *dto.CreateCommentRequest) *SysComment {
	var err error

	// 检查是否已经存在相同的评论
	var existingComment model.Comment
	tx := e.Orm.Debug().Begin()
	err = db.Db.Where("moment_id = ? AND user_id = ? AND contents = ? ", c.MomentID, c.UserID, c.Contents).First(&existingComment).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
	}
	err = db.Db.Create(&model.Comment{
		MomentID: c.MomentID,
		Contents: c.Contents,
		UserID:   c.UserID,
	}).Error
	if err != nil {
		tx.Rollback()

	}

	//更新Moment的Comments数量
	var moment model.Moment
	err = db.Db.First(&moment, c.MomentID).Error
	if err != nil {
		tx.Rollback()
	}
	moment.Comments++

	db.Db.Save(&moment)
	return e
}
