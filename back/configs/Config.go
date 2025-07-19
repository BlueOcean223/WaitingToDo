package configs

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	MySQLConfig    MySQLConfig    `yaml:"mysql"`
	RedisConfig    RedisConfig    `yaml:"redis"`
	MailConfig     MailConfig     `yaml:"mail"`
	MinioConfig    MinioConfig    `yaml:"minio"`
	RabbitMQConfig RabbitMQConfig `yaml:"rabbitmq"`
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
