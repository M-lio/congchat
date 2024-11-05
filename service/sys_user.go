package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
)

type SysUser struct {
	Service
}

func (e *SysUser) GetUser(d *dto.GetUserRequest) *SysUser {
	var err error
	var user model.User
	db.Db.First(&user, d.UserID)

	if user.ID == 0 {
		_ = e.AddError(err)
	}
	return e
}

func (e *SysUser) GetFriends(d *dto.GetFriendsRequest) *SysUser {
	var err error
	var friendships []model.Friendship
	result := db.Db.Preload("User").Preload("Friend").Where("user_id = ? OR friend_id = ?", d.UserID, d.UserID).Find(&friendships)
	if result.Error != nil {
		_ = e.AddError(err)
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
	//更新数据库中的用户信息
	if err = db.Db.Model(&model.User{}).Where("id = ?", d.UserID).Updates(user).Error; err != nil {
		_ = e.AddError(err)
	}
	return e
}
