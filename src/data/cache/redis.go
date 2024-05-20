package cache

import (
	"base_structure/src/config"
	"base_structure/src/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

var redisClient *redis.Client
var redisInit sync.Once
var logger = logging.NewLogger(config.GetConfig())

func InitRedis(cfg *config.Config) error {
	var err error
	redisInit.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:               fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password:           cfg.Redis.Password,
			DB:                 cfg.Redis.Db,
			DialTimeout:        cfg.Redis.DialTimeout * time.Second,
			ReadTimeout:        cfg.Redis.ReadTimeout * time.Second,
			WriteTimeout:       cfg.Redis.WriteTimeout * time.Second,
			PoolSize:           cfg.Redis.PoolSize,
			PoolTimeout:        cfg.Redis.PoolTimeout * time.Second,
			IdleTimeout:        cfg.Redis.IdleTimeout * time.Millisecond,
			IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Millisecond,
		})
		_, err = redisClient.Ping().Result()
		if err != nil {
			redisClient = nil
			logger.Error(logging.Redis, logging.StartUp, "error initializing Redis", nil)
			return
		}
		logger.Info(logging.Redis, logging.StartUp, "Redis connection established", nil)
	})
	return err
}

func GetRedis() *redis.Client {
	if redisClient == nil {
		cfg := config.GetConfig()
		err := InitRedis(cfg)
		if err != nil {
			logger.Fatal(logging.Redis, logging.StartUp, err.Error(), nil)
		}
	}
	return redisClient
}

func CloseRedis() {
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			logger.Fatal(logging.Redis, logging.Closing, "error on closing redis connection", nil)
			return
		}
		redisClient = nil
	}
}

func Set[T any](c *redis.Client, key string, value T, duration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, v, duration).Err()
}

func Get[T any](c *redis.Client, key string) (T, error) {
	dest := *new(T)
	v, err := c.Get(key).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}
