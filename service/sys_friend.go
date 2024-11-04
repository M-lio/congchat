package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
	"gorm.io/gorm"
)

type SysFriends struct {
	Service
}

func (e *SysFriends) AddFriend(c *dto.AddFriendRequest) *SysFriends {
	var err error
	// 检查是否已经存在添加过相同的好友
	var existingFriend model.Friendship
	tx := e.Orm.Debug().Begin()
	err = db.Db.Where("user_id = ? AND friend_username ?", c.UserID, c.FriendUsername).First(&existingFriend).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = e.AddError(err)
		tx.Rollback()
	}

	userID := c.UserID

	var friend model.User

	result := db.Db.First(&friend, "username=?", c.FriendUsername)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}

	var existingFriendship model.Friendship
	result = db.Db.First(&existingFriendship, "user_id = ? AND friend_id = ?", userID, friend.ID)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}

	newFriendship := model.Friendship{
		UserID:   userID,
		FriendID: friend.ID,
		Status:   "pending",
	}

	// 尝试将好友保存到数据库
	result = db.Db.Create(&newFriendship)
	if result.Error != nil {
		_ = e.AddError(err)
		tx.Rollback()
	}
	tx.Commit()
	return e
}

func (e *SysFriends) SearchFriend(c *dto.SearchFriendRequest) *SysFriends {
	var err error
	// 初始化一个用户切片来存储搜索结果
	var users []model.User
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
	if err = db.Db.Where("username = ? OR nick_name = ?OR user_id = ?", c.UserName, c.NickName, c.UserID).Find(&users).Error; err != nil {
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

func (e *SysFriends) AcceptFriend(c *dto.AcceptFriendRequest) *SysFriends {
	var err error
	tx := e.Orm.Debug().Begin()

	// 获取当前用户的 ID
	userID := c.UserID

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

func (e *SysFriends) RejectFriend(c *dto.RejectFriendRequest) *SysFriends {
	var err error
	tx := e.Orm.Debug().Begin()

	// 获取当前用户的 ID
	userID := c.UserID

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
