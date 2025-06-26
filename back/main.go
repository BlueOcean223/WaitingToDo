package main

import (
	"back/configs"
	"back/routers"
	"back/service"
	"back/utils/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:7070", "http://192.168.163.129:7070", "http://101.34.246.32:7070"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization", "Origin", "New-Access-Token"},
		ExposeHeaders: []string{"New-Access-Token"},
	}))

	// 使用JWT令牌校验、拦截
	r.Use(jwt.JWTAuthMiddleware())

	// 加载配置文件
	err := configs.InitConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	// 如果有本地配置文件，则覆盖本地配置
	if _, err := os.Stat("config.local.yaml"); err == nil {
		err = configs.InitConfig("config.local.yaml")
		if err != nil {
			log.Fatalf("加载本地配置文件失败: %v", err)
		}
	}

	// 连接mysql
	err = configs.InitMysqlConnection()
	if err != nil {
		log.Fatalf("数据库连接异常: %v", err)
	}

	// 连接Redis
	err = configs.InitRedisClient()
	if err != nil {
		log.Fatalf("Redis连接异常: %v", err)
	}

	// 连接Minio
	err = configs.InitMinioClient()
	if err != nil {
		log.Fatalf("Minio连接异常: %v", err)
	}

	// 连接RabbitMQ
	err = configs.InitRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ连接异常: %v", err)
	}

	// 启动MQ消费者
	go service.StartFriendConsumer()
	go service.StartTeamConsumer()
	// 启动定时任务
	go service.TickerNotify()

	// 初始化路由
	routers.InitializeRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("程序启动失败")
	}
}
