package cache

import (
	"context"
	"fmt"
	"log"
	"time"
	"woodcarving-backend/config"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis
func InitRedis() {
	cfg := config.GlobalConfig.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	log.Println("Redis初始化成功")
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Delete 删除缓存
func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func Exists(key string) bool {
	result, _ := RedisClient.Exists(ctx, key).Result()
	return result > 0
}
