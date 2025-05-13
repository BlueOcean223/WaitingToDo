package routers

import (
	"back/controllers"
	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		// 登录
		auth.GET("/login", controllers.Login)
	}
}
