package routers

import "github.com/gin-gonic/gin"

// InitializeRoutes 初始化路由
func InitializeRoutes(r *gin.Engine) {
	SetAuthRoutes(r)
	SetUserRoutes(r)
	SetLoadRoutes(r)
	SetTaskRoutes(r)
	SetFriendRoutes(r)
}
