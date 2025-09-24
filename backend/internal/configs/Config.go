package configs

import (
	"backend/pkg/logger"
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

// LogAnalyzerConfig 日志分析器配置
type LogAnalyzerConfig struct {
	Model         string `yaml:"model"`          // 模型名称
	URL           string `yaml:"url"`            // API地址
	APIKey        string `yaml:"api_key"`        // API密钥
	MaxTokens     int    `yaml:"max_tokens"`     // 最大令牌数
	Timeout       int    `yaml:"timeout"`        // 超时时间(秒)
	Enabled       bool   `yaml:"enabled"`        // 是否启用
	AnalysisLevel string `yaml:"analysis_level"` // 分析级别: detailed, quick, trend, custom
	ReportPath    string `yaml:"report_path"`    // 报告存放路径
	AnalysisHours int    `yaml:"analysis_hours"` // 分析时间范围(小时)
	MaxLogSize    int    `yaml:"max_log_size"`   // 最大日志大小(字节)
}

type Config struct {
	MySQLConfig       MySQLConfig       `yaml:"mysql"`
	RedisConfig       RedisConfig       `yaml:"redis"`
	MailConfig        MailConfig        `yaml:"mail"`
	MinioConfig       MinioConfig       `yaml:"minio"`
	RabbitMQConfig    RabbitMQConfig    `yaml:"rabbitmq"`
	LogConfig         LogConfig         `yaml:"log"`
	LogAnalyzerConfig LogAnalyzerConfig `yaml:"log_analyzer"`
}

var AppConfigs Config

func initConfig(configPath string) error {
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

func InitConfig() {
	// 加载配置文件
	err := initConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	// 如果有本地配置文件，则覆盖
	if _, err := os.Stat("./config/config.local.yaml"); err == nil {
		err = initConfig("./config/config.local.yaml")
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
