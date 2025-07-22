package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	Client *redis.Client
}

func NewRedisLimiter(addr, password string, db int) *RedisLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisLimiter{Client: client}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string, limit, blockDuration int) (bool, error) {
	now := time.Now().Unix()
	k := fmt.Sprintf("rl:%s:%d", key, now)

	count, err := r.Client.Incr(ctx, k).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.Client.Expire(ctx, k, time.Second)
	}

	if count > int64(limit) {
		blockKey := fmt.Sprintf("block:%s", key)
		r.Client.Set(ctx, blockKey, "1", time.Duration(blockDuration)*time.Second)
		return false, nil
	}

	block, _ := r.Client.Get(ctx, fmt.Sprintf("block:%s", key)).Result()
	return block == "", nil
}
