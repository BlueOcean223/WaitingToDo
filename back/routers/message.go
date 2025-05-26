package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetMessageRoutes(r *gin.Engine) {
	// 初始化依赖
	messageRepository := repository.NewMessageRepository(configs.MysqlDb)
	messageService := service.NewMessageService(messageRepository)
	messageController := controllers.NewMessageController(messageService)

	message := r.Group("/message")
	{
		// 获取用户未读消息数量
		message.GET("/unreadCount", messageController.GetUnreadMessageCount)
		// 分页查询消息
		message.GET("/list", messageController.GetMessageList)
		// 更新消息
		message.PUT("/update", messageController.UpdateMessage)
		// 删除消息
		message.DELETE("/delete", messageController.DeleteMessage)
		// 一键已读
		message.PUT("/readAll", messageController.ReadAllMessage)
		// 处理请求
		message.POST("/handle", messageController.HandleRequest)
	}
}
