package main

import (
	"backend/internal/routers"
	"backend/internal/services/consumer"
	"backend/internal/services/ticker"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 启动MQ消费者
	startConsumers()
	// 启动定时任务
	go ticker.TaskNotify()

	// 初始化路由
	routers.InitializeRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("程序启动失败")
	}
}

// 启动MQ消费者
func startConsumers() {
	// 创建消费者管理器
	consumerManager := consumer.NewConsumerManager()

	// 添加消费者
	consumerManager.RegisterConsumer(consumer.NewFriendConsumer())
	consumerManager.RegisterConsumer(consumer.NewTeamConsumer())

	// 启动所有消费者服务
	consumerManager.StartAll()

	// 优雅关闭
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		log.Println("收到关闭信号，正在退出")

		consumerManager.Stop()
		os.Exit(0)
	}()
}
