package handlers

import (
	"backend/internal/models/dto"
	"backend/internal/models/vo"
	"backend/internal/services"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FriendHandler struct {
	friendService services.FriendService
}

func NewFriendHandler(friendService services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

// GetFriendInfo 根据id查询好友详情
func (s *FriendHandler) GetFriendInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	friendInfo, err := s.friendService.GetFriendInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}
	if friendInfo == (dto.FriendInfoDto{}) {
		c.JSON(http.StatusOK, response.Fail("", "用户不存在", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", friendInfo))
}

// GetFriendList 获取好友列表
func (s *FriendHandler) GetFriendList(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	friendList, err := s.friendService.GetFriendList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", friendList))
}

// SearchUserByEmail 根据邮箱搜索用户
func (s *FriendHandler) SearchUserByEmail(c *gin.Context) {
	searchEmail := c.Query("email")
	// 获取用户邮箱
	userEmail := c.GetString("user")

	user, err := s.friendService.SearchUserByEmail(userEmail, searchEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}
	if user == (dto.FriendInfoDto{}) {
		c.JSON(http.StatusOK, response.Success("", "用户不存在", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", user))
}

// AddFriend 发送添加好友请求
func (s *FriendHandler) AddFriend(c *gin.Context) {
	var friendVo vo.FriendVo
	if err := c.ShouldBind(&friendVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
	}

	err := s.friendService.AddFriend(friendVo.UserId, friendVo.FriendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}
	c.JSON(http.StatusOK, response.Success("", "发送成功", nil))
}

// DeleteFriend  删除好友
func (s *FriendHandler) DeleteFriend(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("userId"))
	friendId, err := strconv.Atoi(c.Query("friendId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err = s.friendService.DeleteFriend(userId, friendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "操作成功", nil))
	return
}
