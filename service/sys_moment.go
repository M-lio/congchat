package service

import "C"
import (
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"fmt"
)

type SysMoment struct {
	Service
}

// CreateMoment 创建SysComment
func (e *SysMoment) CreateMoment(d *dto.CreateMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()
	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in CreateMoment:", r)
		}
	}()

	// 文本参数校验 查看是否符合格式
	if err = dto.ValidateContent(d.Content); err != nil {
		err = errors.New("内容验证错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//检验下图片的输入是否大于9张图片的处理函数
	if err = dto.ValidateImgURLs(d.ImgURL); err != nil {
		err = errors.New("图片验证错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 创建新的Moment并保存到数据库（校验）
	moment := &model.Moment{
		UserID:  d.UserID,
		Content: d.Content,
		ImgURL:  d.ImgURL,
	}

	result := tx.Create(moment)
	if result.Error != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()
	return e
}

func (e *SysMoment) EditMoment(d *dto.EditMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()
	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in EditMoment:", r)
		}
	}()

	// 文本参数校验 查看是否符合格式
	if err = dto.ValidateContent(d.Content); err != nil {
		err = errors.New("内容验证错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//检验下图片的输入是否大于9张图片的处理函数
	if err = dto.ValidateImgURLs(d.ImgURL); err != nil {
		err = errors.New("验证错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//初始化一个新的结构体来接受
	var originalMoment model.Moment

	result := tx.First(&originalMoment, d.ID)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	originalMoment.Content = d.Content
	originalMoment.ImgURL = d.ImgURL

	// 保存更改到数据库
	if err = tx.Save(&originalMoment).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	return e
}

func (e *SysMoment) GetMoment(c *dto.GetMomentRequest, list *[]model.Moment) *SysMoment {
	var err error

	tx := e.Orm.Debug().Begin()
	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in GetMoment:", r)
		}
	}()

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

	query := tx.Model(&model.Moment{}).Preload("Comments").Order("created_at DESC")
	if len(whereConditions) > 0 {
		query = query.Where(whereConditions, whereArgs...)
	}
	err = query.Offset(offset).Limit(c.Limit).Find(&list).Error

	// 处理查询结果
	if err != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	tx.Commit()

	return e
}

// RemoveMoment 删除SysComment
func (e *SysMoment) RemoveMoment(d *dto.DeleteMomentRequest) *SysMoment {
	var err error
	tx := e.Orm.Debug().Begin()

	//参数校验（需要通过才能继续往下走）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in RemoveMoment:", r)
		}
	}()

	// 参数校验
	if d.Ids == nil || len(d.Ids) == 0 {
		err = errors.New("必须提供至少一个Moment ID")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	for _, id := range d.Ids {
		if id <= 0 {
			err = fmt.Errorf("无效的Moment ID: %d", id)
			e.handleErrorAndRollback(tx, err)
			return e
		}
	}
	//初始化moment动态 根据momentID来删除动态
	var moment model.Moment
	//绑参
	result := tx.Where("id = ?", d.Ids).Delete(&moment)
	//校验
	if result.Error != nil {
		err = result.Error
		e.handleErrorAndRollback(tx, err)
		return e
	}
	if result.RowsAffected == 0 {
		err = errors.New("未找到要删除的Moment记录")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()

	return e
}
