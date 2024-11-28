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

	// 调用 AddFriendClear 函数检查输入
	if err = dto.AddFriendClear(d); err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
	}

	// 检查是否已经存在添加过相同的好友关系
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

	// 将好友关系添加到 Redis 缓存中
	key := fmt.Sprintf("friendship:%d:%d", d.UserID, friend.ID)
	_, err = db.RedisClient.SAdd(ctx, key, "pending").Result() // 存储一个状态
	if err != nil {
		tx.Rollback()
		e.AddError(errors.New("添加好友关系到 Redis 缓存时发生错误"))
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

	// 调用 SearchFriendClear 函数检查输入
	if err = dto.SearchFriendClear(d); err != nil {
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
	if err = tx.Where("username = ? OR nick_name = ?OR user_id = ?", d.UserName, d.NickName, d.UserID).Find(&users).Error; err != nil {
		e.handleErrorAndRollback(tx, err)
		return e
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
	tx.Commit()
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

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, d.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ = e.AddError(err)
		}
		tx.Rollback()
	}

	// 检查当前用户是否是接受该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != d.UserID && friendship.FriendID != d.UserID {
		_ = e.AddError(err)
		tx.Rollback()
	}

	// 更新好友请求的状态为“已接受”
	db.Db.Model(&friendship).Update("Status", "accepted")
	if db.Db.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	// 更新 Redis 中的好友关系状态
	key := fmt.Sprintf("friendship:%d:%d", friendship.UserID, friendship.FriendID)
	// 使用 HSet 或 HDel 根据需要添加或删除关系，这里假设我们只需要更新状态
	_, err = db.RedisClient.HSet(ctx, key, "status", "accepted").Result()
	if err != nil {
		tx.Rollback()
		_ = e.AddError(errors.New("更新好友关系到 Redis 缓存时发生错误"))
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()
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

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, d.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ = e.AddError(err)
		}
		tx.Rollback()
	}

	// 检查当前用户是否是拒绝该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != d.UserID && friendship.FriendID != d.UserID {
		_ = e.AddError(err)
		tx.Rollback()
	}

	// 更新好友请求的状态为“已拒绝”
	db.Db.Model(&friendship).Update("Status", "rejected")
	if db.Db.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	// 更新 Redis 中的好友关系状态
	key := fmt.Sprintf("friendship:%d:%d", friendship.UserID, friendship.FriendID)
	// 使用 HSet 或 HDel 根据需要添加或删除关系，这里假设我们只需要更新状态
	_, err = db.RedisClient.HSet(ctx, key, "status", "rejected").Result()
	if err != nil {
		tx.Rollback()
		_ = e.AddError(errors.New("更新好友关系到 Redis 缓存时发生错误"))
		e.handleErrorAndRollback(tx, err)
		return e
	}
	tx.Commit()
	return e
}
