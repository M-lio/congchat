package service

import "C"
import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
)

type SysMoment struct {
	Service
}

// CreateMoment 创建SysComment
func (e *SysMoment) CreateMoment(d *dto.CreateMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()
	userID := d.UserID
	//TODO 参数校验（是否符合需求，万一涉及暴力，黄色，广告。直接拦截）

	// 创建新的Moment并保存到数据库（校验）
	moment := &model.Moment{
		UserID:  userID,
		Content: d.Content,
		ImgURL:  d.ImgURL,
		Goods:   d.Goods,
		GoodsID: d.GoodsID,
	}
	result := db.Db.Create(moment)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	tx.Commit()
	return e
}

func (e *SysMoment) EditMoment(d *dto.EditMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()
	//初始化一个新的结构体来接受
	var originalMoment model.Moment

	result := db.Db.First(&originalMoment, d.ID)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	originalMoment.Content = d.Content
	originalMoment.ImgURL = d.ImgURL

	// 保存更改到数据库
	if err = db.Db.Save(&originalMoment).Error; err != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}

	return e
}

func (e *SysMoment) GetMoment(c *dto.GetMomentRequest, list *[]model.Moment) *SysMoment {
	var err error

	// 构建查询（校验）旧代码 这里我通过传一个进来
	//var moments []model.Moment

	// 限制Limit的最大值为10
	if c.Limit > 10 {
		c.Limit = 10
	}

	// 预处理分页参数
	offset := c.StartIdx

	// 初始化查询条件
	var whereConditions []interface{}
	var whereArgs []interface{}
	if c.FriendID != 0 {
		whereConditions = append(whereConditions, "user_id = ?")
		whereArgs = append(whereArgs, c.FriendID)
	}
	if c.SearchStr != "" {
		searchStr := "%" + c.SearchStr + "%"
		whereConditions = append(whereConditions, "content LIKE ?")
		whereArgs = append(whereArgs, searchStr)
	}

	query := db.Db.Model(&model.Moment{}).Preload("Comments").Order("created_at DESC")
	if len(whereConditions) > 0 {
		query = query.Where(whereConditions, whereArgs...)
	}
	err = query.Offset(offset).Limit(c.Limit).Find(&list).Error

	// 处理查询结果
	if err != nil {
		_ = e.AddError(err)

	}
	return e
}

// RemoveComment 删除SysComment
func (e *SysMoment) RemoveMoment(d *dto.DeleteMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()
	//初始化moment动态 根据momentID来删除动态
	var moment model.Moment
	//绑参
	result := db.Db.Where("id = ?", d.Ids).Delete(&moment)
	//校验
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	if result.RowsAffected == 0 {
		_ = e.AddError(err)
		tx.Rollback()
	}
	tx.Commit()

	return e
}
