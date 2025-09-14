package handlers

import (
	"backend/internal/models"
	"backend/internal/models/vo"
	"backend/internal/services"
	"backend/pkg/myError"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// GetUnreadMessageCount 获取用户未读消息数量
func (s *MessageHandler) GetUnreadMessageCount(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	count, err := s.messageService.GetUnreadMessageCount(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", count))
}

// GetMessageList 获取用户消息列表
func (s *MessageHandler) GetMessageList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	messageList, err := s.messageService.GetMessageList(page, pageSize, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", messageList))
}

// UpdateMessage 更新消息
func (s *MessageHandler) UpdateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}
	err := s.messageService.UpdateMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "更新成功", nil))
}

// DeleteMessage 删除消息
func (s *MessageHandler) DeleteMessage(c *gin.Context) {
	messageId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err = s.messageService.DeleteMessage(messageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "删除成功", nil))
}

// ReadAllMessage 一键已读
func (s *MessageHandler) ReadAllMessage(c *gin.Context) {
	var userId int
	if err := c.ShouldBindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.ReadAllMessage(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "操作成功", nil))
}

// HandleRequest 同意请求
func (s *MessageHandler) HandleRequest(c *gin.Context) {
	var messageVo vo.MessageVo
	if err := c.ShouldBindJSON(&messageVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.HandleRequest(messageVo)
	if err != nil {
		if myError.IsMyError(err) {
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, response.Success("", "操作成功", nil))
}

// AddMessage 添加消息
func (s *MessageHandler) AddMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.messageService.AddMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "添加成功", nil))
}
