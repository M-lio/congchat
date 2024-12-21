package service

import (
	"fmt"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Msg   string
	MsgID string
	Error error
}

func (db *Service) AddError(err error) error {
	if db.Error == nil {
		db.Error = err
	} else if err != nil {
		db.Error = fmt.Errorf("%v; %w", db.Error, err)
	}
	return db.Error
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysComment) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysFriends) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysUser) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysMoment) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysGoods) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}

// 辅助函数，用于处理错误并回滚事务
func (e *SysTransf) handleErrorAndRollback(tx *gorm.DB, err error) {
	fmt.Println("错误:", err)
	_ = e.AddError(err) // 假设 e 有一个 AddError 方法来记录错误
	tx.Rollback()
}
