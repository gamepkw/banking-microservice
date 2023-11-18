package repository_test

import (
	"context"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockTransactionRepo struct {
	mockDB    *sqlmock.Sqlmock
	mockRedis *MockRedisClient
}

func NewMockTransactionRepository(mockDB sqlmock.Sqlmock, mockRedis *MockRedisClient) MockTransactionRepository {
	return &mockTransactionRepo{
		mockDB:    &mockDB,
		mockRedis: mockRedis,
	}
}

type MockTransactionRepository interface {
	GetTransactionConfig(ctx context.Context, configName string) (string, error)
}

func (m *mockTransactionRepo) GetTransactionConfig(ctx context.Context, configName string) (string, error) {
	fmt.Println("Get value from db")
	configValue := "100"

	return configValue, nil
}
