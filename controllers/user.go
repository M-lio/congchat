package controllers

import (
	"congchat-user/db"
	"congchat-user/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Friendship struct {
	gorm.Model
	model.User
	UserID   uint   `gorm:"not null"`
	FriendID uint   `gorm:"not null"`
	Status   string `gorm:"default:'pending'"`
}
type FriendshipStatus struct {
	FriendID uint   `json:"friend_id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

// 控制器包含需要用到的实体 以及逻辑处理函数
// 获取用户资料//修改时间10.24.1.1
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

// 更新用户资料//修改时间10.24.1.2
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

func GetFriendsHandler(c *gin.Context) {
	userID := c.GetInt("id")

	var friendships []Friendship
	result := db.Db.Preload("User").Preload("Friend").Where("user_id = ? OR friend_id = ?", userID, userID).Find(&friendships)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch friendships"})
		return
	}

	var friendStatuses []FriendshipStatus //为什么用切片形式
	for _, friendship := range friendships {
		if friendship.UserID == uint(userID) {
			var friend model.User
			db.Db.First(&friend, friendship.FriendID)
			friendStatuses = append(friendStatuses, FriendshipStatus{
				FriendID: friendship.FriendID,
				Username: friend.Username,
				Status:   friendship.Status,
			})
		} else {
			var user model.User
			db.Db.First(&user, friendship.UserID)
			friendStatuses = append(friendStatuses, FriendshipStatus{
				FriendID: friendship.UserID,
				Username: user.Username,
				Status:   friendship.Status,
			})
		}
	}

	c.JSON(http.StatusOK, friendStatuses)
}
