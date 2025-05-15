package routers

import (
	"back/controllers"
	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		// 登录
		auth.POST("/login", controllers.Login)
		// 注册
		auth.POST("/register", controllers.Register)
		// 忘记密码
		auth.POST("/forget", controllers.Forget)
		// 发送验证码
		auth.GET("/captcha", controllers.Captcha)
	}
}
