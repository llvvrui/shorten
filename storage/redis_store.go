package storage

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/llvvrui/shortener/utils"
	"github.com/redis/go-redis/v9"
)

var ErrShortCodeExists = errors.New("short code already exists")

type Store interface {
	SaveShortUrl(shortCode, longUrl string, expiration time.Duration) (string, error)
	GetLongUrl(shortCode string) (string, error)
	IncrementVisit(shortCode string) error
	GetVisitCount(shortCode string) (int64, error)
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) SaveShortUrl(shortCode, longUrl string, expiration time.Duration) (string, error) {
	ctx := context.Background()
	var generatedCode string

	if shortCode == "" {
		id, err := r.client.Incr(ctx, "counter").Result()
		if err != nil {
			return "", err
		}
		generatedCode = utils.Base62Encode(id)
		shortCode = generatedCode
	} else {
		key := "short:" + shortCode
		result, err := r.client.SetArgs(ctx, key, longUrl, redis.SetArgs{
			Mode: "NX",
			TTL:  expiration,
		}).Result()
		if err != nil {
			return "", err
		}
		if result == "" {
			return "", ErrShortCodeExists
		}
		statsKey := "stats:" + shortCode
		if err := r.client.Set(ctx, statsKey, 0, expiration).Err(); err != nil {
			r.client.Del(ctx, key)
			return "", err
		}
		return shortCode, nil
	}

	key := "short:" + shortCode
	if err := r.client.Set(ctx, key, longUrl, expiration).Err(); err != nil {
		return "", err
	}
	statsKey := "stats:" + shortCode
	if err := r.client.Set(ctx, statsKey, 0, expiration).Err(); err != nil {
		return "", err
	}
	return shortCode, nil
}

func (r *RedisStore) GetLongUrl(shortCode string) (string, error) {
	ctx := context.Background()
	key := "short:" + shortCode
	longUrl, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("short code not found")
	}
	return longUrl, err
}

func (r *RedisStore) IncrementVisit(shortCode string) error {
	ctx := context.Background()
	statsKey := "stats:" + shortCode
	_, err := r.client.Incr(ctx, statsKey).Result()
	return err
}

func (r *RedisStore) GetVisitCount(shortCode string) (int64, error) {
	ctx := context.Background()
	statsKey := "stats:" + shortCode
	countStr, err := r.client.Get(ctx, statsKey).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(countStr, 10, 64)
}

