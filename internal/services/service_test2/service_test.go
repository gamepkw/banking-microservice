package service

import (
	"time"

	repo "github.com/gamepkw/transactions-banking-microservice/internal/repositories/repositories_test"
	"github.com/stretchr/testify/mock"
)

type mockTransactionService struct {
	mock.Mock
	mockTransactionRepo repo.MockTransactionRepository
	contextTimeout      time.Duration
	mockRedis           *MockRedisClient
}

func NewMockTransactionService(
	tr repo.MockTransactionRepository,
	timeout time.Duration,
	redis *MockRedisClient,
) *mockTransactionService {
	return &mockTransactionService{
		mockTransactionRepo: tr,
		contextTimeout:      timeout,
		mockRedis:           redis,
	}
}
