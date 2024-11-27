package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"context"
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
	// 调用 CreateClear 函数检查输入
	if err = dto.CreateClear(d); err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	err = tx.Where("moment_id = ? AND user_id = ? AND contents = ? ", d.MomentID, d.UserID, d.Contents).First(&existingComment).Error
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		err = errors.New("已存在相同的评论")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 创建新的评论并保存到数据库
	comment := &model.Comment{
		MomentID: d.MomentID,
		Contents: d.Contents,
		UserID:   d.UserID,
	}
	err = tx.Create(comment).Error
	if err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 更新Redis中的评论列表和评论数量（使用Set存储评论ID）
	momentKey := fmt.Sprintf("moment:%d", d.MomentID)
	commentSetKey := fmt.Sprintf("%s:comment_ids", momentKey)

	// 将评论ID添加到Redis Set中
	_, err = db.RedisClient.SAdd(context.Background(), commentSetKey, comment.ID).Result()
	if err != nil {
		tx.Rollback()
		fmt.Println("Error adding comment ID to Redis set:", err)
		return e
	}

	// 更新评论数量（这里可以优化为只使用Redis来维护评论数量，但为了确保数据一致性，还是先从数据库中获取当前数量）
	var moment model.Moment
	err = tx.First(&moment, d.MomentID).Error
	if err != nil {
		tx.Rollback()
		e.handleErrorAndRollback(tx, err)
		return e
	}
	moment.Comments++

	// 将更新后的评论数量缓存到Redis中（可以设置过期时间）
	commentCountKey := fmt.Sprintf("%s:comments_count", momentKey)
	_, err = db.RedisClient.Set(context.Background(), commentCountKey, moment.Comments, 0).Result() // 0表示没有设置过期时间，根据需求设置
	if err != nil {
		tx.Rollback()
		// 处理Redis错误（在实际应用中可能需要记录日志或采取其他措施）
		fmt.Println("Error setting comment count in Redis:", err)
		return e
	}

	// 提交事务
	tx.Commit()
	return e
}

// RemoveComment 删除SysComment
func (e *SysComment) RemoveComment(d *dto.DeleteCommentRequest) *SysComment {
	var err error
	var data model.Comment
	tx := e.Orm.Debug().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in RemoveComment:", r)
		}
	}()

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

	// 从Redis Set中移除评论ID
	momentKey := fmt.Sprintf("moment:%d", data.MomentID)
	commentSetKey := fmt.Sprintf("%s:comment_ids", momentKey)
	_, err = db.RedisClient.SRem(context.Background(), commentSetKey, d.Ids).Result()
	if err != nil {
		tx.Rollback()
		return e
	}

	// 更新Redis中的评论数量（如果需要的话）
	// 这里可以优化为直接从Redis中减少数量，但为了数据一致性，这里仍然从数据库中获取并更新
	var moment model.Moment
	err = tx.First(&moment, data.MomentID).Error
	if err != nil {
		tx.Rollback()
		return e
	}
	moment.Comments--

	commentCountKey := fmt.Sprintf("%s:comments_count", momentKey)
	_, err = db.RedisClient.Set(context.Background(), commentCountKey, moment.Comments, 0).Result() // 0表示没有设置过期时间
	if err != nil {
		tx.Rollback()
		return e
	}
	tx.Commit()

	return e
}
