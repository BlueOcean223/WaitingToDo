package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetFriendRoutes(r *gin.Engine) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	friendRepository := repository.NewFriendRepository(configs.MysqlDb)
	messageRepository := repository.NewMessageRepository(configs.MysqlDb)

	friendService := service.NewFriendService(authRepository, friendRepository, messageRepository)

	friendController := controllers.NewFriendController(friendService)

	// 配置路由
	friend := r.Group("/friend")
	{
		// 根据id查询好友详情
		friend.GET("/info", friendController.GetFriendInfo)
		// 获取好友列表
		friend.GET("/list", friendController.GetFriendList)
		// 根据邮箱查询用户信息，用于添加好友
		friend.GET("/search", friendController.SearchUserByEmail)
		// 添加好友请求
		friend.POST("/add", friendController.AddFriend)
		// 删除好友
		friend.DELETE("/delete", friendController.DeleteFriend)
	}
}
