package service

import (
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	result := db.Db.First(&user, d.UserID)
	if result.Error != nil {
		return e // 返回一个空的 SysUser 实例
	}

	userKey := fmt.Sprintf("user:%d", d.UserID)

	// 将用户信息序列化并存储到 Redis 中
	userData, err := json.Marshal(user)
	if err != nil {
		// 序列化失败，返回“空”的 SysUser 实例
		// 注意：这里可以记录日志或执行其他错误处理逻辑
		return e
	}

	// 设置 Redis 键值对，没有设置过期时间
	_, err = db.RedisClient.Set(context.Background(), userKey, userData, 0).Result()
	if err != nil {
		// Redis 存储失败，返回“空”的 SysUser 实例
		// 注意：这里可以记录日志、执行重试逻辑或向监控系统报告错误
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

	userKey := fmt.Sprintf("user:%d", d.UserID)
	userData, err := json.Marshal(user)
	if err != nil {
		err = errors.New("序列化用户数据时发生错误")
		return e
	}
	_, err = db.RedisClient.Set(context.Background(), userKey, userData, 0).Result() // 0表示没有设置过期时间
	if err != nil {
		err = errors.New("更新Redis用户信息时发生错误")
		return e
	}
	return e
}
