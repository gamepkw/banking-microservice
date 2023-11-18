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

// NewMockTransactionService creates a new instance of mockTransactionService with dependencies
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

type Transaction struct {
	Id          int64     `json:"id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Fee         float64   `json:"fee"`
	Total       float64   `json:"total"`
	SubmittedAt time.Time `json:"submitted_at"`
	CreatedAt   time.Time `json:"created_at"`
	Account     Account   `json:"account"`
	Receiver    Account   `json:"receiver,omitempty"`
}

type Account struct {
	AccountNo string     `json:"account_no,omitempty"`
	Uuid      string     `json:"uuid,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Tel       string     `json:"tel,omitempty"`
	Balance   float64    `json:"balance"`
	Bank      string     `json:"bank,omitempty"`
	Status    string     `json:"status,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	IsClosed  int        `json:"is_closed,omitempty"`
}

// type mockRedis interface {
// 	Set(key string, value interface{}, exp time.Duration) error
// 	Get(key string) (string, error)
// }

// // repository represent the repository model
// type mockredis struct {
// 	Client redis.Cmdable
// }

// // NewRedisRepository will create an object that represent the Repository interface
// func NewRedisRepository(Client redis.Cmdable) mockRedis {
// 	return &mockredis{Client}
// }

// // Set attaches the redis repository and set the data
// func (r *mockredis) Set(key string, value interface{}, exp time.Duration) error {
// 	return r.Client.Set(key, value, exp).Err()
// }

// // Get attaches the redis repository and get the data
// func (r *mockredis) Get(key string) (string, error) {
// 	get := r.Client.Get(key)
// 	return get.Result()
// }
