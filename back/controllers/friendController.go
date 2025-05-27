package controllers

import (
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FriendController struct {
	friendService *service.FriendService
}

func NewFriendController(friendService *service.FriendService) *FriendController {
	return &FriendController{friendService: friendService}
}

// GetFriendInfo 根据id查询好友详情
func (s *FriendController) GetFriendInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	friendInfo, err := s.friendService.GetFriendInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}
	if friendInfo == (dto.FriendInfoDto{}) {
		c.JSON(http.StatusOK, models.Fail("", "用户不存在", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", friendInfo))
}

// GetFriendList 获取好友列表
func (s *FriendController) GetFriendList(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	friendList, err := s.friendService.GetFriendList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", friendList))
}

// SearchUserByEmail 根据邮箱搜索用户
func (s *FriendController) SearchUserByEmail(c *gin.Context) {
	searchEmail := c.Query("email")
	// 获取用户邮箱
	userEmail, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}

	user, err := s.friendService.SearchUserByEmail(userEmail, searchEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}
	if user == (dto.FriendInfoDto{}) {
		c.JSON(http.StatusOK, models.Success("", "用户不存在", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", user))
}

// AddFriend 发送添加好友请求
func (s *FriendController) AddFriend(c *gin.Context) {
	var friendVo vo.FriendVo
	if err := c.ShouldBind(&friendVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
	}

	err := s.friendService.AddFriend(friendVo.UserId, friendVo.FriendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}
	c.JSON(http.StatusOK, models.Success("", "发送成功", nil))
}

// DeleteFriend  删除好友
func (s *FriendController) DeleteFriend(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("userId"))
	friendId, err := strconv.Atoi(c.Query("friendId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	err = s.friendService.DeleteFriend(userId, friendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "操作成功", nil))
	return
}
