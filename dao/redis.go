package dao

import (
	"Voichatter/configs"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

func InitRedis(cfg *configs.RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort), // 地址和端口
		Username: cfg.RedisUsername,                                  // 用户名
		Password: cfg.RedisPassword,                                  // 设置密码
		DB:       cfg.RedisDbName,                                    // DB名称
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		fmt.Println("Redis 未启动或初始化失败")
		panic(err)
	}
	RedisClient = client
}
