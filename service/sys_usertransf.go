package service

import (
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"fmt"
	"time"
)

type SysTransf struct {
	Service
}

func (e *SysTransf) CreateTransfer(d *dto.TransferRequest) *SysTransf {
	var err error
	tx := e.Orm.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in CreateTransfer:", r)
			// 在实际应用中，你可能希望将错误记录到日志中，而不是仅仅打印到控制台
		}
	}()
	// 参数校验
	if d.SenderID == 0 || d.ReceiverID == 0 {
		err = errors.New("sender_id and receiver_id are required")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	if d.Amount <= 0 {
		err = errors.New("amount must be greater than zero")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 更新发送转账用户的余额
	var sender model.User
	if err := tx.First(&sender, d.SenderID).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}
	//判断下余额够不够
	if sender.Balance < d.Amount {
		err = errors.New("insufficient balance")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	//更新当前发送转账用户余额
	sender.Balance -= d.Amount
	if err := tx.Save(&sender).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 创建转账记录
	transfer := model.Transfer{
		SenderID:   d.SenderID,
		ReceiverID: d.ReceiverID,
		Amount:     d.Amount,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	if err := tx.Create(&transfer).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		return e
	}

	// 开启一个grotine来监听
	go func() {
		time.Sleep(24 * time.Hour)
		tx := e.Orm.Begin()
		defer tx.Rollback() // 注意：这里的Rollback是为了防止在查询或更新时出错，如果操作成功则不应调用

		var pendingTransfers []model.Transfer
		if err := tx.Where("status = ? AND created_at <= ?", "pending", time.Now().Add(-24*time.Hour)).Find(&pendingTransfers).Error; err != nil {
			fmt.Println("Error querying pending transfers:", err)
			return
		}

		for _, tf := range pendingTransfers {
			var receiver model.User
			if err := tx.First(&receiver, tf.ReceiverID).Error; err != nil {
				fmt.Println("Error finding receiver:", err)
				continue
			}
			receiver.Balance += tf.Amount
			if err := tx.Save(&receiver).Error; err != nil {
				fmt.Println("Error updating receiver balance:", err)
				continue
			}

			tf.Status = "refunded"
			if err := tx.Save(&tf).Error; err != nil {
				fmt.Println("Error updating transfer status:", err)
				continue
			}

			var sender model.User
			if err := tx.First(&sender, tf.SenderID).Error; err != nil {
				fmt.Println("Error finding sender:", err)
				continue
			}
			sender.Balance += tf.Amount
			if err := tx.Save(&sender).Error; err != nil {
				fmt.Println("Error updating sender balance:", err)
				continue
			}
		}

		if err := tx.Commit().Error; err != nil {
			fmt.Println("Error committing refund transaction:", err)
		}
	}()

	return e
}
