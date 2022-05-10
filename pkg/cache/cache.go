package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"order/pkg/configs"
	"sync"
	"time"
)

const (
	redisKeyFormat = "%v:%v"
	prefix         = "tinder"
)

type ProcessData func(data interface{}) error

type Adapter interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}, expiration time.Duration) error
	SetWithFunc(key string, v interface{}, expiration time.Duration, validate ProcessData) error
}

type RedisCache struct {
	ctx    context.Context
	client *redis.Client
	lock   sync.Mutex
	prefix string
}

var adapter Adapter

func InitRedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     configs.RedisURL(),
		Password: configs.RedisPassword(),
		DB:       0,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic(fmt.Sprintf("Cannot connect redis client: %s", err.Error()))
	}
	adapter = &RedisCache{
		client: client,
		prefix: prefix,
		lock:   sync.Mutex{},
	}
}

func Instance() Adapter {
	return adapter
}

func (rc *RedisCache) Get(key string, v interface{}) error {
	originalKey := fmt.Sprintf(redisKeyFormat, rc.prefix, key)
	rc.lock.Lock()
	defer rc.lock.Unlock()
	data, err := rc.client.Get(originalKey).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (rc *RedisCache) Set(key string, v interface{}, expiration time.Duration) error {
	originalKey := fmt.Sprintf(redisKeyFormat, rc.prefix, key)
	rc.lock.Lock()
	defer rc.lock.Unlock()
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	_, err = rc.client.Set(originalKey, data, expiration).Result()
	return err
}

func (rc *RedisCache) SetWithFunc(key string, v interface{}, expiration time.Duration, process ProcessData) error {
	originalKey := fmt.Sprintf(redisKeyFormat, rc.prefix, key)
	rc.lock.Lock()
	defer rc.lock.Unlock()
	data, err := rc.client.Get(originalKey).Bytes()
	if err := process(data); err != nil {
		return err
	}
	jData, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	_, err = rc.client.Set(originalKey, jData, expiration).Result()
	return err
}
