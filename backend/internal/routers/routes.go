package routers

import (
	"backend/internal/middlewares/jwt"
	"backend/internal/middlewares/security"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitializeRoutes 初始化路由
func InitializeRoutes(r *gin.Engine) {
	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:5173", "http://localhost:7070", "http://192.168.163.129:7070", "http://101.34.246.32:7070"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization", "Origin", "New-Access-Token"},
		ExposeHeaders: []string{"New-Access-Token"},
	}))

	// 注册全局中间件
	r.Use(security.SecurityMiddleware())
	r.Use(security.RateLimitMiddleware())

	// 放行的路由
	SetAuthRoutes(r)

	// 需要校验的路由
	// 使用JWT令牌校验、拦截
	protected := r.Group("/")
	protected.Use(jwt.JWTAuthMiddleware())
	{
		SetUserRoutes(protected)
		SetLoadRoutes(protected)
		SetTaskRoutes(protected)
		SetFriendRoutes(protected)
		SetMessageRoutes(protected)
	}
}
