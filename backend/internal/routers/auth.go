package routers

import (
	"backend/internal/configs"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(r *gin.Engine) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	authService := services.NewAuthService(authRepository)
	authController := handlers.NewAuthHandler(authService)

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
