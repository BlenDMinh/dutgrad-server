package services

import (
    "context"
    "encoding/json"
    "time"

    "github.com/BlenDMinh/dutgrad-server/configs"
    "github.com/redis/go-redis/v9"
)

type RedisService struct {
    client *redis.Client
}

func NewRedisService() *RedisService {
    config := configs.GetEnv().Redis
    client := redis.NewClient(&redis.Options{
        Addr:     config.Addr,
        Password: config.Password,
        DB:       config.DB,
    })

    return &RedisService{client: client}
}

func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
    json, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return s.client.Set(context.Background(), key, json, expiration).Err()
}

func (s *RedisService) Get(key string) (string, error) {
    return s.client.Get(context.Background(), key).Result()
}

func (s *RedisService) Del(key string) error {
    return s.client.Del(context.Background(), key).Err()
}
