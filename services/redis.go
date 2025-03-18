package services

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/BlenDMinh/dutgrad-server/configs"
    "github.com/redis/go-redis/v9"
)

type RedisService struct {
    client *redis.Client
}

var instance *RedisService = nil

func NewRedisService() *RedisService {
    if instance != nil {
        return instance
    }
    config := configs.GetEnv().Redis
    log.Printf("Connecting to Redis at %s...", config.Addr)
    
    client := redis.NewClient(&redis.Options{
        Addr:     config.Addr,
        Password: config.Password,
        DB:       config.DB,
    })

    // Test the connection
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        log.Printf("Failed to connect to Redis: %v", err)
        return &RedisService{client: client}
    }
    
    log.Printf("Successfully connected to Redis at %s", config.Addr)
    instance = &RedisService{client: client}
    return instance
}

func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
    json, err := json.Marshal(value)
    if err != nil {
        log.Printf("Redis marshal error for key %s: %v", key, err)
        return err
    }
    
    err = s.client.Set(context.Background(), key, json, expiration).Err()
    if err != nil {
        log.Printf("Redis SET error for key %s: %v", key, err)
        return err
    }
    
    log.Printf("Successfully set Redis key: %s (expires in %v)", key, expiration)
    return nil
}

func (s *RedisService) Get(key string) (string, error) {
    if s.client == nil {
        log.Printf("Redis client is nil")
        return "", nil
    }
    val, err := s.client.Get(context.Background(), key).Result()
    if err != nil {
        if err == redis.Nil {
            log.Printf("Redis key not found: %s", key)
        } else {
            log.Printf("Redis GET error for key %s: %v", key, err)
        }
        return "", err
    }
    
    log.Printf("Successfully retrieved Redis key: %s", key)
    return val, nil
}

func (s *RedisService) Del(key string) error {
    err := s.client.Del(context.Background(), key).Err()
    if err != nil {
        log.Printf("Redis DEL error for key %s: %v", key, err)
        return err
    }
    
    log.Printf("Successfully deleted Redis key: %s", key)
    return nil
}
