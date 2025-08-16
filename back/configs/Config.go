package configs

import (
	"back/utils/logger"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	Filename   string `yaml:"filename"`    // 日志文件名
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留的日志文件数量
	MaxAge     int    `yaml:"max_age"`     // 日志文件保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
	Console    bool   `yaml:"console"`     // 是否输出到控制台
}

type Config struct {
	MySQLConfig    MySQLConfig    `yaml:"mysql"`
	RedisConfig    RedisConfig    `yaml:"redis"`
	MailConfig     MailConfig     `yaml:"mail"`
	MinioConfig    MinioConfig    `yaml:"minio"`
	RabbitMQConfig RabbitMQConfig `yaml:"rabbitmq"`
	LogConfig      LogConfig      `yaml:"log"`
}

var AppConfigs Config

func InitConfig(configPath string) error {
	// 读取yaml文件
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 解析yaml文件到结构体
	err = yaml.Unmarshal(yamlFile, &AppConfigs)
	if err != nil {
		return err
	}

	return nil
}

// InitLogger 初始化日志系统
func InitLogger() error {
	return logger.InitLogger(logger.LogConfig(AppConfigs.LogConfig))
}

func init() {
	// 加载配置文件
	err := InitConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	// 如果有本地配置文件，则覆盖
	if _, err := os.Stat("config.local.yaml"); err == nil {
		err = InitConfig("config.local.yaml")
		if err != nil {
			log.Fatalf("加载本地配置文件失败: %v", err)
		}
	}

	// 初始化日志系统
	err = InitLogger()
	if err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 连接mysql
	err = InitMysqlConnection()
	if err != nil {
		log.Fatalf("数据库连接异常: %v", err)
	}

	// 连接Redis
	err = InitRedisClient()
	if err != nil {
		log.Fatalf("Redis连接异常: %v", err)
	}

	// 连接Minio
	err = InitMinioClient()
	if err != nil {
		log.Fatalf("Minio连接异常: %v", err)
	}

	// 连接RabbitMQ
	err = InitRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ连接异常: %v", err)
	}
}
