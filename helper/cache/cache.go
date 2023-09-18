package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace-svc/helper/config"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type CacheDatabase interface {
	Deletes(keys ...string)
	Delete(key string)
	DeletePrefix(prefix string)
	Expire(key string, expiration time.Duration)
	Get(key string) (string, error)
	GetJsonStruct(key string, structObj interface{}) error
	HashDelete(key string, field string)
	HashDeletes(key string, fields ...string)
	HashGet(key string, field string) (string, error)
	HashGetJsonStruct(key string, field string, structObj interface{}) error
	IsExist(key string) bool
	Set(key string, value interface{}, exp ...int) error
	SetJsonStruct(key string, value interface{}, exp ...int)
	HashSet(key string, field string, value interface{}, exp ...int) error
	RPush(key string, value interface{})
	LRem(key string, count int64, value interface{})
}

type redisEngine struct {
	RedisClient *redis.Client
	RedisUse    bool
	Expires     time.Duration
}

// SetupRedisConnection setup redis connection.
func SetupRedisConnection(cfg *config.CacheDBConfig) (CacheDatabase, error) {
	var redisClient *redis.Client
	if cfg.IsUse {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.Index,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return nil, err
		}
	}

	return redisEngine{
		RedisClient: redisClient,
		RedisUse:    cfg.IsUse,
		Expires:     time.Duration(int64(cfg.DefaultExpiry)) * time.Minute,
	}, nil
}

func (cache redisEngine) Deletes(keys ...string) {
	if cache.RedisUse {
		for _, key := range keys {
			if cache.IsExist(key) {
				cache.RedisClient.Del(ctx, key)
			}
		}
	}
}

func (cache redisEngine) Delete(key string) {
	if cache.RedisUse {
		if cache.IsExist(key) {
			cache.RedisClient.Del(ctx, key)
		}
	}
}

func (cache redisEngine) HashDelete(key string, field string) {
	if cache.RedisUse {
		if cache.IsExist(key) {
			cache.RedisClient.HDel(ctx, key, field)
		}
	}
}

func (cache redisEngine) HashDeletes(key string, fields ...string) {
	if cache.RedisUse {
		for _, field := range fields {
			if cache.IsExist(key) {
				cache.RedisClient.HDel(ctx, key, field)
			}
		}
	}
}

func (cache redisEngine) DeletePrefix(prefix string) {
	if cache.RedisUse {
		iter := cache.RedisClient.Scan(ctx, 0, fmt.Sprintf("%v*", prefix), 0).Iterator()
		for iter.Next(ctx) {
			cache.RedisClient.Del(ctx, iter.Val())
		}
	}
}

func (cache redisEngine) Expire(key string, expiration time.Duration) {
	if cache.RedisUse && cache.IsExist(key) {
		cache.RedisClient.Expire(ctx, key, expiration)
	}
}

func (cache redisEngine) Get(key string) (string, error) {
	if cache.RedisUse {
		val, err := cache.RedisClient.Get(ctx, key).Result()
		if err != nil {
			return "", err
		}
		return val, nil
	}
	return "", nil
}

func (cache redisEngine) GetJsonStruct(key string, structObj interface{}) error {
	if cache.RedisUse {

		val, err := cache.RedisClient.Get(ctx, key).Bytes()

		if err == redis.Nil {
			return nil
		}

		if err != nil {
			return err
		}

		err = json.Unmarshal(val, &structObj)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (cache redisEngine) HashGet(key string, field string) (string, error) {
	if cache.RedisUse {
		val, err := cache.RedisClient.HGet(ctx, key, field).Result()
		if err != nil {
			return "", err
		}
		return strings.Trim(val, "\""), nil
	}
	return "", nil
}

func (cache redisEngine) HashGetJsonStruct(key string, field string, structObj interface{}) error {
	if cache.RedisUse {
		val, err := cache.RedisClient.HGet(ctx, key, field).Bytes()
		if err != nil {
			return err
		}

		err = json.Unmarshal(val, &structObj)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (cache redisEngine) IsExist(key string) bool {
	if cache.RedisUse {
		_, err := cache.RedisClient.Get(ctx, key).Result()
		return err != redis.Nil
	}
	return cache.RedisUse
}

func (cache redisEngine) Set(key string, value interface{}, exp ...int) error {
	if cache.RedisUse {
		setExp := cache.Expires
		if exp != nil {
			getExp := exp[0]
			setExp = time.Duration(getExp) * time.Minute
		}
		_, err := cache.RedisClient.Set(ctx, key, value, setExp).Result()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (cache redisEngine) HashSet(key string, field string, value interface{}, exp ...int) error {
	if cache.RedisUse {
		setExp := cache.Expires
		if exp != nil {
			getExp := exp[0]
			setExp = time.Duration(getExp) * time.Minute
		}

		b, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("cache: Marshal key=" + key + " failed: " + err.Error())
		}

		cache.RedisClient.Do(ctx, "HSET", key, field, b, "EX", setExp) //HSetNX(ctx, key, value, field, setExp)
		if err != nil {
			v := string(b)
			if len(v) > 15 {
				v = v[0:12] + "..."
			}
			return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
		}

		return nil
	}

	return nil
}

func (cache redisEngine) SetJsonStruct(key string, value interface{}, exp ...int) {
	if cache.RedisUse {
		setExp := cache.Expires

		if exp != nil {
			getExp := exp[0]
			setExp = time.Duration(getExp) * time.Minute
		}
		json, _ := json.Marshal(value)
		cache.RedisClient.Set(ctx, key, string(json), setExp)
	}
}

func (cache redisEngine) RPush(key string, value interface{}) {
	if cache.RedisUse {
		cache.RedisClient.RPush(ctx, key, value)
	}
}

func (cache redisEngine) LRem(key string, count int64, value interface{}) {
	if cache.RedisUse {
		cache.RedisClient.LRem(ctx, key, count, value)
	}
}
