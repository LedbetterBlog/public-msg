package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// RedisPoolManager 结构体
type RedisPoolManager struct {
	client *redis.Client
}

// NewRedisPoolManager 创建 RedisPoolManager 实例
func NewRedisPoolManager(addr, password string, db int) (*RedisPoolManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, // 使用配置文件中的数据库编号
	})

	// 测试连接
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	log.Println("Redis 连接成功")
	return &RedisPoolManager{client: client}, nil
}

// GetUniqueID 从 Redis 获取唯一 ID
func (r *RedisPoolManager) GetUniqueID(ctx context.Context) (string, error) {
	id, err := r.client.Incr(ctx, "my_counter").Result()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// GetValue 从 Redis 获取值
func (r *RedisPoolManager) GetValue(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetValue 向 Redis 设置值
func (r *RedisPoolManager) SetValue(ctx context.Context, key, value string) error {
	_, err := r.client.Set(ctx, key, value, 0).Result()
	return err
}

// SetSAddValue 向 Redis 设置集合
func (r *RedisPoolManager) SetSAddValue(ctx context.Context, key, value string) error {
	_, err := r.client.SAdd(ctx, key, value).Result()
	return err
}

// GetSMEMBERSValue 从 Redis 获取集合的所有成员
func (r *RedisPoolManager) GetSMEMBERSValue(ctx context.Context, key string) ([]string, error) {
	// 使用 SMEMBERS 方法获取集合的所有成员
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

// RemoveFromSet 从 Redis 集合中删除一个或多个值
func (r *RedisPoolManager) RemoveFromSet(ctx context.Context, key string, values ...string) error {
	// 使用 SREM 方法从集合中删除指定的值
	_, err := r.client.SRem(ctx, key, values).Result()
	if err != nil {
		return fmt.Errorf("error removing values from set: %w", err)
	}
	return nil
}
