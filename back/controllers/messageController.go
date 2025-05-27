package controllers

import (
	"back/models"
	"back/models/vo"
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

// ReadAllMessage 一键已读
func (s *MessageController) ReadAllMessage(c *gin.Context) {
	var userId int
	if err := c.ShouldBindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.ReadAllMessage(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "操作成功", nil))
}

// HandleRequest 同意请求
func (s *MessageController) HandleRequest(c *gin.Context) {
	var messageVo vo.MessageVo
	if err := c.ShouldBindJSON(&messageVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.HandleRequest(messageVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, models.Success("", "操作成功", nil))
}

// AddMessage 添加消息
func (s *MessageController) AddMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.AddMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, models.Success("", "添加成功", nil))
}
