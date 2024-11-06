package service

import (
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SysComment struct {
	Service
}

func (e *SysComment) CreateComment(d *dto.CreateCommentRequest) *SysComment {
	var err error
	tx := e.Orm.Debug().Begin()
	// 参数校验 检查是否已经存在相同的评论
	var existingComment model.Comment

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in CreateComment:", r)
		}
	}()
	// 文本参数校验 查看是否符合格式
	if err = dto.ValidateContents(d.Contents); err != nil {
		err = errors.New("内容验证格式错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 其他参数校验，例如检查MomentID是否为0（假设0是无效的朋友圈ID）
	if d.MomentID == 0 {
		err = errors.New("朋友圈ID不能为空")
		err = errors.New("朋友圈ID错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	err = tx.Where("moment_id = ? AND user_id = ? AND contents = ? ", d.MomentID, d.UserID, d.Contents).First(&existingComment).Error
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	err = tx.Create(&model.Comment{
		MomentID: d.MomentID,
		Contents: d.Contents,
		UserID:   d.UserID,
	}).Error
	if err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//更新Moment的Comments数量
	var moment model.Moment
	err = tx.First(&moment, d.MomentID).Error
	if err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}
	moment.Comments++

	tx.Save(&moment)
	tx.Commit()
	return e
}

// RemoveComment 删除SysComment
func (e *SysComment) RemoveComment(d *dto.DeleteCommentRequest) *SysComment {
	var err error
	var data model.Comment
	tx := e.Orm.Debug().Begin()

	// 参数校验
	if len(d.Ids) == 0 {
		err = errors.New("ID 列表不能为空")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	tx.First(&data, d.Ids)
	if err = tx.First(&data, d.Ids).Error; err != nil {
		err = errors.New("并未找到该ID的记录")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	if err = tx.Delete(&data, d.Ids).Error; err != nil {
		err = errors.New("删除该ID中的记录时发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 更新 Moment 的 Comments 数量（）
	var moment model.Moment
	if err = tx.First(&moment, data.MomentID).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}
	moment.Comments--
	if err = tx.Save(&moment).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()

	return e
}
