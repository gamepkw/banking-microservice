package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type MockRedisClient struct {
	Data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	mockRedis := &MockRedisClient{
		Data: make(map[string]string),
	}
	mockRedis.Set("min_deposit_amount", "100", -1)
	// mockRedis.Set("min_deposit_amount", "100", -1)

	return mockRedis
}

// func NewMockRedisClient() *MockRedisClient {
// 	return &MockRedisClient{
// 		Data: make(map[string]string),
// 	}
// }

func (m *MockRedisClient) Set(key string, value interface{}, exp time.Duration) error {
	m.Data[key] = value.(string)
	return nil
}

func (m *MockRedisClient) Get(key string) (string, error) {
	if value, ok := m.Data[key]; ok {
		fmt.Println("Get value from cache")
		return value, nil
	}
	return "", redis.Nil
}

type RedisClient interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
