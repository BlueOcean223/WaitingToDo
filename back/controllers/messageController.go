package controllers

import (
	"back/models"
	"back/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageController struct {
	messageService *service.MessageService
}

func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

// GetUnreadMessageCount 获取用户未读消息数量
func (s *MessageController) GetUnreadMessageCount(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	count, err := s.messageService.GetUnreadMessageCount(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", count))
}

// GetMessageList 获取用户消息列表
func (s *MessageController) GetMessageList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	messageList, err := s.messageService.GetMessageList(page, pageSize, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", messageList))
}

// UpdateMessage 更新消息
func (s *MessageController) UpdateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}
	err := s.messageService.UpdateMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "更新成功", nil))
}

// DeleteMessage 删除消息
func (s *MessageController) DeleteMessage(c *gin.Context) {
	messageId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	err = s.messageService.DeleteMessage(messageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "删除成功", nil))
}
