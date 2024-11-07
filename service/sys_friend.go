package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SysFriends struct {
	Service
}

func (e *SysFriends) AddFriend(d *dto.AddFriendRequest) *SysFriends {
	var err error
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in AddFriend:", r)
		}
	}()

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID不能为空")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//检查FriendUsername是否为空
	if d.FriendUsername == "" {
		err = errors.New("好友用户名不能为空")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 检查是否已经存在添加过相同的好友
	var existingFriend model.Friendship
	err = tx.Where("user_id = ? AND friend_username ?", d.UserID, d.FriendUsername).First(&existingFriend).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("没有找到用户ID或者名字并发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	userID := d.UserID

	var friend model.User

	result := tx.First(&friend, "username=?", d.FriendUsername)
	if result.Error != nil {
		err = errors.New("查找好友用户时发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	var existingFriendship model.Friendship
	result = tx.First(&existingFriendship, "user_id = ? AND friend_id = ?", userID, friend.ID)
	if result.Error != nil {
		err = errors.New("查询好友记录时发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	newFriendship := model.Friendship{
		UserID:   userID,
		FriendID: friend.ID,
		Status:   "pending",
	}

	// 尝试将好友保存到数据库
	result = tx.Create(&newFriendship)
	if result.Error != nil {
		err = errors.New("添加好友时发生错误")
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()
	return e
}

func (e *SysFriends) SearchFriend(d *dto.SearchFriendRequest) *SysFriends {
	var err error
	// 初始化一个用户切片来存储搜索结果
	var users []model.User
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in SearchFriend:", r)
		}
	}()

	//检查UserName是否为空
	if d.UserName == "" {
		err = errors.New("用户名不能为空")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	//检查UserID是否为空
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	/*/ 从查询参数中获取搜索关键词
	//searchKey := c.Query("search")
	//if searchKey == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Search key is required"})
	//	return
	//}
	*/
	//这里我改成通过请求中的用户id或者昵称或者用户名来查询用户，最终返回
	// 在数据库中搜索用户
	// 注意：这里使用了OR查询来同时搜索用户名和手机号
	if err = db.Db.Where("username = ? OR nick_name = ?OR user_id = ?", d.UserName, d.NickName, d.UserID).Find(&users).Error; err != nil {
		_ = e.AddError(err)
	}
	// 如果找到匹配的用户，则返回用户列表
	// 注意：这里返回了用户的部分信息，您可以根据需要调整返回的字段
	userList := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			// 可以添加其他需要返回的字段，比如头像、个性签名等
		}
	}

	return e
}

func (e *SysFriends) AcceptFriend(d *dto.AcceptFriendRequest) *SysFriends {
	var err error
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in AcceptFriend:", r)
		}
	}()

	//检查UserID是否为空
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 获取当前用户的 ID
	userID := d.UserID

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ = e.AddError(err)
		}
		tx.Rollback()
	}

	// 检查当前用户是否是接受该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != userID && friendship.FriendID != userID {
		_ = e.AddError(err)
		tx.Rollback()
	}

	// 更新好友请求的状态为“已接受”
	db.Db.Model(&friendship).Update("Status", "accepted")
	if db.Db.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	return e
}

func (e *SysFriends) RejectFriend(d *dto.RejectFriendRequest) *SysFriends {
	var err error
	tx := e.Orm.Debug().Begin()

	//defer语句来捕获可能发生的恐慌（panic），并在发生恐慌时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in RejectFriend:", r)
		}
	}()

	//检查UserID是否为空
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 获取当前用户的 ID
	userID := d.UserID

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ = e.AddError(err)
		}
		tx.Rollback()
	}

	// 检查当前用户是否是拒绝该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != userID && friendship.FriendID != userID {
		_ = e.AddError(err)
		tx.Rollback()
	}

	// 更新好友请求的状态为“已拒绝”
	db.Db.Model(&friendship).Update("Status", "rejected")
	if db.Db.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	return e
}
