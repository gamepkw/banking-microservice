package repository_test

import (
	"errors"
)

// MockRedisClient is a custom mock implementation of a Redis client.
type MockRedisClient struct {
	Data map[string]string
}

// NewMockRedisClient creates a new instance of the mock Redis client.
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Data: make(map[string]string),
	}
}

// Set sets a key-value pair in the mock Redis.
func (m *MockRedisClient) Set(key string, value string) error {
	m.Data[key] = value
	return nil
}

// Get retrieves a value by key from the mock Redis.
func (m *MockRedisClient) Get(key string) (string, error) {
	if value, ok := m.Data[key]; ok {
		return value, nil
	}
	return "", errors.New("Key not found")
}

type RedisClient interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
