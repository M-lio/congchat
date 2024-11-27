package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"context"
	"errors"
	"fmt"
	"log"
)

type SysGoods struct {
	Service
}

var ctx = context.Background()

func Rgoods(d *dto.AddGoodRequest) {
	key := fmt.Sprintf("moment:%s:goods", d.MomentID)
	// 使用 SAdd 命令将用户ID添加到集合中
	_, err := db.RedisClient.SAdd(ctx, key, d.UserID).Result()
	if err != nil {
		log.Fatalf("Failed to add like: %v", err)
	}
}

func RremoveGoods(d *dto.CancelGoodRequest) {
	key := fmt.Sprintf("moment:%s:goods", d.MomentID)
	// 使用 SRem 命令从集合中移除用户ID
	_, err := db.RedisClient.SRem(ctx, key, d.UserID).Result()
	if err != nil {
		log.Fatalf("Failed to remove good: %v", err)
	}
}

func (e *SysGoods) AddGood(d *dto.AddGoodRequest) *SysGoods {
	var err error
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in SearchFriend:", r)
		}
	}()

	// 调用 IsClear 函数检查输入
	if err = dto.AisClear(d); err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 检查是否已经点赞过
	var Good model.Goods
	result := db.Db.First(&Good, "user_id = ? AND moment_id = ?", d.UserID, d.MomentID)
	if result.Error != nil {
		err = errors.New("无法找到朋友圈记录")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	if result.RowsAffected == 0 { // RowsAffected 返回受影响的行数，如果没有找到则为0
		// 创建新的点赞记录
		newGood := model.Goods{UserID: d.UserID, MomentID: d.MomentID}
		if err := db.Db.Create(&newGood).Error; err != nil {
			err = errors.New("创建新的点赞记录时发生错误")
			e.handleErrorAndRollback(tx, err)
			return e
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
			err = errors.New("更新点赞记录时发生错误")
			_ = e.AddError(err)
			tx.Rollback()
		}

	}
	Rgoods(&dto.AddGoodRequest{d.MomentID, d.UserID})
	tx.Commit()

	return e
}

func (e *SysGoods) CancelGood(d *dto.CancelGoodRequest) *SysGoods {
	var err error
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in SearchFriend:", r)
		}
	}()

	// 调用 IsClear 函数检查输入
	if err = dto.CisClear(d); err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 查找并删除点赞记录
	var Good model.Goods // 确保 Like 结构体中有 UserID 和 MomentID 字段
	result := db.Db.Delete(&Good, "user_id = ? AND moment_id = ?", d.UserID, d.MomentID)
	if result.Error != nil {
		// 处理数据库删除错误
		err = errors.New("删除点赞时发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	if result.RowsAffected == 0 {
		// 如果没有找到要删除的记录，则返回错误或提示信息
		err = errors.New("没有找到要删除的记录发生错误")
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
			err = errors.New("更新点赞时发生错误")
			_ = e.AddError(err)
			tx.Rollback()
		}
	}
	RremoveGoods(&dto.CancelGoodRequest{d.UserID, d.MomentID})
	tx.Commit()

	return e
}
