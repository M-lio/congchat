package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
)

type SysGoods struct {
	Service
}

func (e *SysGoods) AddGood(d *dto.AddGoodRequest) *SysGoods {
	var err error
	tx := e.Orm.Debug().Begin()
	// 检查是否已经点赞过
	var Good model.Goods
	result := db.Db.First(&Good, "user_id = ? AND moment_id = ?", d.UserID, d.MomentID)
	if result.Error != nil {
		// 处理数据库查询错误
		_ = e.AddError(err)
		tx.Rollback()
	}
	if result.RowsAffected == 0 { // RowsAffected 返回受影响的行数，如果没有找到则为0
		// 创建新的点赞记录
		newGood := model.Goods{UserID: d.UserID, MomentID: d.MomentID}
		if err := db.Db.Create(&newGood).Error; err != nil {
			// 处理创建记录错误
			_ = e.AddError(err)
			tx.Rollback()
		}

		// 更新朋友圈动态的点赞数
		var moment model.Moment
		if err := db.Db.First(&moment, d.MomentID).Error; err != nil {
			// 处理查询动态错误
			_ = e.AddError(err)
			tx.Rollback()
		}
		moment.Goods++
		if err := db.Db.Save(&moment).Error; err != nil {
			// 处理更新点赞数错误
			_ = e.AddError(err)
			tx.Rollback()
		}
	}
	tx.Commit()

	return e
}

func (e *SysGoods) CancelGood(d *dto.CancelGoodRequest) *SysGoods {
	var err error
	tx := e.Orm.Debug().Begin()
	// 查找并删除点赞记录
	var Good model.Goods // 确保 Like 结构体中有 UserID 和 MomentID 字段
	result := db.Db.Delete(&Good, "user_id = ? AND moment_id = ?", d.UserID, d.MomentID)
	if result.Error != nil {
		// 处理数据库删除错误
		_ = e.AddError(err)
		tx.Rollback()
	}
	if result.RowsAffected == 0 {
		// 如果没有找到要删除的记录，则返回错误或提示信息
		_ = e.AddError(err)
		tx.Rollback()
	}

	// 更新朋友圈动态的点赞数
	var moment model.Moment
	if err := db.Db.First(&moment, d.MomentID).Error; err != nil {
		// 处理查询动态错误
		_ = e.AddError(err)
		tx.Rollback()
	}
	// 假设 Goods 字段存储了点赞数，这里需要减一
	if moment.Goods > 0 {
		moment.Goods--
		if err := db.Db.Save(&moment).Error; err != nil {
			// 处理更新点赞数错误
			_ = e.AddError(err)
			tx.Rollback()
		}
	}
	tx.Commit()

	return e
}
