package routers

import (
	"back/controllers"
	"github.com/gin-gonic/gin"
)

func SetUserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	{
		// 发送验证码
		user.POST("/checkCaptcha", controllers.CheckCaptcha)
		// 重置密码
		user.POST("/reset", controllers.Reset)
		// 修改个人信息
		user.POST("/update", controllers.UpdateUserInfo)
	}
}
