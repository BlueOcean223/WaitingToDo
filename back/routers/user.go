package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetUserRoutes(r *gin.RouterGroup) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	authService := service.NewAuthService(authRepository)
	userService := service.NewUserService(authRepository)
	userController := controllers.NewUserController(authService, userService)

	user := r.Group("/user")
	{
		// 发送验证码
		user.POST("/checkCaptcha", userController.CheckCaptcha)
		// 重置密码
		user.POST("/reset", userController.Reset)
		// 修改个人信息
		user.POST("/update", userController.UpdateUserInfo)
		// 获取用户信息
		user.GET("/info", userController.GetUserInfo)
	}
}
