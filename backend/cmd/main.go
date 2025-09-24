package main

import (
	"backend/internal/configs"
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

	// 配置代理信任设置
	var trustedProxies = []string{
		"127.0.0.1",      // 本地回环
		"::1",            // IPv6 本地回环
		"172.16.0.0/12",  // Docker 默认网络段
		"10.0.0.0/8",     // 私有网络段
		"192.168.0.0/16", // 私有网络段
	}
	err := r.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("设置代理信任失败:", err)
	}

	// 初始化配置
	configs.InitConfig()

	// 初始化路由
	routers.InitializeRoutes(r)

	// 启动MQ消费者
	consumerManager := startConsumers()

	// 启动定时任务
	ticker.StartAllTimers()

	// 设置优雅关闭
	go setupGracefulShutdown(consumerManager)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("程序启动失败")
	}
}

// 启动MQ消费者
func startConsumers() *consumer.ConsumerManager {
	// 创建消费者管理器
	consumerManager := consumer.NewConsumerManager()

	// 添加消费者
	consumerManager.RegisterConsumer(consumer.NewFriendConsumer())
	consumerManager.RegisterConsumer(consumer.NewTeamConsumer())

	// 启动所有消费者服务
	consumerManager.StartAll()

	return consumerManager
}

// 设置优雅关闭
func setupGracefulShutdown(consumerManager *consumer.ConsumerManager) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	log.Println("收到关闭信号，正在优雅关闭服务...")

	// 停止定时任务
	ticker.StopAllTimers()

	// 停止MQ消费者
	consumerManager.Stop()

	log.Println("所有服务已停止")
	os.Exit(0)
}
