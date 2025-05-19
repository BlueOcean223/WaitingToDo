package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(r *gin.Engine) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	authService := service.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	auth := r.Group("/auth")
	{
		// 登录
		auth.POST("/login", authController.Login)
		// 注册
		auth.POST("/register", authController.Register)
		// 忘记密码
		auth.POST("/forget", authController.Forget)
		// 发送验证码
		auth.GET("/captcha", authController.Captcha)
	}
}
