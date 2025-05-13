package configs

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	MySQLConfig MySQLConfig `yaml:"mysql"`
	RedisConfig RedisConfig `yaml:"redis"`
	MailConfig  MailConfig  `yaml:"mail"`
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
