package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"errors"
)

type SysUser struct {
	Service
}

func (e *SysUser) GetUser(d *dto.GetUserRequest) *SysUser {
	var err error
	var user model.User

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		return e
	}

	db.Db.First(&user, d.UserID)

	if user.ID == 0 {
		err = errors.New("用户ID不能为0")
		e.handleErrorAndRollback(db.Db, err)
		return e
	}

	return e
}

func (e *SysUser) GetFriends(d *dto.GetFriendsRequest) *SysUser {
	var err error
	var friendships []model.Friendship

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		return e
	}

	if err = db.Db.Preload("User").Preload("Friend").Where("user_id = ? OR friend_id = ?", d.UserID, d.UserID).Find(&friendships).Error; err != nil {
		err = errors.New("查询好友关系无法找到该记录")
		return e
	}

	var friendStatuses []model.FriendshipStatus //为什么用切片形式
	for _, friendship := range friendships {
		if friendship.UserID == d.UserID {
			var friend model.User
			db.Db.First(&friend, friendship.FriendID)
			friendStatuses = append(friendStatuses, model.FriendshipStatus{
				FriendID: friendship.FriendID,
				Username: friend.Username,
				Status:   friendship.Status,
			})
		} else {
			var user model.User
			db.Db.First(&user, friendship.UserID)
			friendStatuses = append(friendStatuses, model.FriendshipStatus{
				FriendID: friendship.UserID,
				Username: user.Username,
				Status:   friendship.Status,
			})
		}
	}
	return e
}

func (e *SysUser) UpdateUser(d *dto.UpdateUserRequest) *SysUser {
	var err error
	var user model.User

	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		err = errors.New("用户ID不能为0")
		return e
	}
	//更新数据库中的用户信息
	if err = db.Db.Model(&model.User{}).Where("id = ?", d.UserID).Updates(user).Error; err != nil {
		err = errors.New("更新用户信息时发生错误")
		return e
	}
	return e
}
