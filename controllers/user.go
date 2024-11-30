package controllers

import (
	"congchat-user/core"
	"congchat-user/model"
	"congchat-user/service"
	"congchat-user/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SysUser struct {
	core.Api
}

// Get 处理获取用户资料的接口
// 修改时间11.6
func (e SysUser) Get(c *gin.Context) {
	s := service.SysUser{}
	var rsp core.Rsp
	req := dto.GetUserRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := s.GetUser(&req).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp.Code = 0
	rsp.Data = user
	rsp.Msg = "获取用户资料成功"
	c.JSON(http.StatusOK, rsp)
}

// 控制器包含需要用到的实体 以及逻辑处理函数 旧代码
// 获取用户资料//修改时间10.24.1.1
/*
func GetUserHandle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user model.User
	db.Db.First(&user, id)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return
	}

	c.JSON(http.StatusOK, user)
}
*/

// Update 更新用户资料
func (e SysUser) Update(c *gin.Context) {
	s := service.SysUser{}
	var rsp core.Rsp
	req := dto.UpdateUserRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.UpdateUser(&req).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp.Code = 0
	rsp.Msg = "更新用户资料成功"
	c.JSON(http.StatusOK, rsp)
}

// 更新用户资料//修改时间10.24.1.2 旧代码
/*
func UpdateUserHandle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user model.User

	//绑定请求体到结构体中
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//更新数据库中的用户信息
	db.Db.Model(&model.User{}).Where("id = ?", id).Updates(user)
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
*/

// GetFriends 获取好友列表
func (e SysUser) GetFriends(c *gin.Context) {
	s := service.SysUser{}
	var rsp core.Rsp
	req := dto.GetFriendsRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var FList model.FriendshipStatusList
	err := s.GetFriends(&req).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp.Code = 0
	rsp.Data = FList
	rsp.Msg = "获取好友列表成功"
	c.JSON(http.StatusOK, rsp)
}

/*
func GetFriendsHandler(c *gin.Context) {
	userID := c.GetInt("id")

	var friendships []model.Friendship
	result := db.Db.Preload("User").Preload("Friend").Where("user_id = ? OR friend_id = ?", userID, userID).Find(&friendships)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch friendships"})
		return
	}

	var friendStatuses []model.FriendshipStatus //为什么用切片形式
	for _, friendship := range friendships {
		if friendship.UserID == uint(userID) {
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

	c.JSON(http.StatusOK, friendStatuses)
}
*///GetFriends 获取好友列表旧代码
