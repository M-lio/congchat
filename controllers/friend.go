package controllers

import (
	"congchat-user/core"
	"congchat-user/db"
	"congchat-user/model"
	"congchat-user/service"
	"congchat-user/service/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type SysFriends struct {
	core.Api
}

// 添加好友接口2
func (e SysFriends) Add(c *gin.Context) {
	var req dto.AddFriendRequest
	var rsp core.Rsp
	s := new(service.SysFriends)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.AddFriend(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Msg = "添加好友成功"
	c.JSON(http.StatusOK, rsp)
}

/*// 添加好友接口2
func AddFriendHandler(c *gin.Context) {
	var req dto.AddFriendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("id")

	var friend model.User

	result := db.Db.First(&friend, "username=?", req.FriendUsername)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingFriendship model.Friendship
	result = db.Db.First(&existingFriendship, "user_id = ? AND friend_id = ?", userID, friend.ID)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Friendship already exists"})
		return
	}

	newFriendship := model.Friendship{
		UserID:   uint(userID),
		FriendID: friend.ID,
		Status:   "pending",
	}

	// 尝试将好友保存到数据库
	result = db.Db.Create(&newFriendship)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create friendship"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}
*/ //添加好友接口旧版本代码

// 搜索好友的处理器函数接口1
func (e SysFriends) Search(c *gin.Context) {
	var req dto.SearchFriendRequest
	var rsp core.Rsp
	s := new(service.SysFriends)
	// 初始化一个用户切片来存储搜索结果
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回具体的验证错误信息
		return
	}
	err := s.SearchFriend(&req).Error
	if err != nil {
		return
	}

	rsp.Code = 0
	rsp.Data = req.UserName
	rsp.Msg = "搜素好友成功"
	c.JSON(http.StatusOK, rsp)
}

/*// 搜索好友的处理器函数接口1
func SearchFriendsHandler(c *gin.Context) {
	// 从查询参数中获取搜索关键词
	searchKey := c.Query("search")
	if searchKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search key is required"})
		return
	}

	// 初始化一个用户切片来存储搜索结果
	var users []model.User

	// 在数据库中搜索用户
	// 注意：这里使用了OR查询来同时搜索用户名和手机号
	result := db.Db.Where("username = ? OR phone = ?", searchKey, searchKey).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
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

	c.JSON(http.StatusOK, gin.H{"users": userList})
}
*/

// acceptFriendsRequestHandler 接受好友申请的处理器3
func AcceptFriendsRequestHandler(c *gin.Context) {
	// 从 URL 参数中获取好友请求 ID
	friendshipIDStr := c.Param("id")
	friendshipID, err := strconv.ParseUint(friendshipIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friendship ID"})
		return
	}

	// 获取当前用户的 ID（这里假设已经通过中间件设置了）
	userID := c.GetInt64("user_id")

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, uint(friendshipID))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friendship not found"})
		}
		return
	}

	// 检查当前用户是否是接受该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != uint(userID) && friendship.FriendID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to accept this friendship"})
		return
	}

	// 更新好友请求的状态为“已接受”
	db.Db.Model(&friendship).Update("Status", "accepted")
	if db.Db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friendship status"})
		return
	}

	// 可选：执行其他逻辑，如更新双方用户的好友列表、发送通知等

	c.JSON(http.StatusOK, gin.H{"message": "Friendship accepted"})
}

func RejectFriendsRequestHandler(c *gin.Context) {
	// 从 URL 参数中获取好友请求 ID
	friendshipIDStr := c.Param("id")
	friendshipID, err := strconv.ParseUint(friendshipIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friendship ID"})
		return
	}

	// 获取当前用户的 ID（这里假设已经通过中间件设置了）
	userID := c.GetInt64("user_id")

	var friendship model.Friendship
	// 查找好友请求
	result := db.Db.First(&friendship, uint(friendshipID))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friendship not found"})
		}
		return
	}

	// 检查当前用户是否是拒绝该请求的合适人选（例如，检查请求是否发送给当前用户）
	if friendship.UserID != uint(userID) && friendship.FriendID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to reject this friendship"})
		return
	}

	// 更新好友请求的状态为“已拒绝”
	db.Db.Model(&friendship).Update("Status", "rejected")
	if db.Db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friendship status"})
		return
	}

	// 可选：执行其他逻辑，如发送通知给申请者等

	c.JSON(http.StatusOK, gin.H{"message": "Friendship rejected"})
}
