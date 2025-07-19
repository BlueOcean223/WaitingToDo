package main

import (
	"back/routers"
	"back/service"
	"back/utils/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:5173", "http://localhost:7070", "http://192.168.163.129:7070", "http://101.34.246.32:7070"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization", "Origin", "New-Access-Token"},
		ExposeHeaders: []string{"New-Access-Token"},
	}))

	// 使用JWT令牌校验、拦截
	r.Use(jwt.JWTAuthMiddleware())

	// 启动MQ消费者
	go service.StartFriendConsumer()
	go service.StartTeamConsumer()
	// 启动定时任务
	go service.TickerNotify()

	// 初始化路由
	routers.InitializeRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("程序启动失败")
	}
}
