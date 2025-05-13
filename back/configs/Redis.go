package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

var RedisClient *redis.Client

func InitRedisClient() error {
	redisConfig := AppConfigs.RedisConfig
	// 连接redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})

	// 检查连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
