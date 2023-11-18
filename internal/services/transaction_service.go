package service

import (
	"context"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	repo "github.com/gamepkw/transactions-banking-microservice/internal/repositories"

	"github.com/go-redis/redis"
)

type transactionService struct {
	transactionRepo repo.TransactionRepository
	contextTimeout  time.Duration
	redis           *redis.Client
}

func NewTransactionService(
	tr repo.TransactionRepository,
	timeout time.Duration,
	redis *redis.Client,
) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		contextTimeout:  timeout,
		redis:           redis,
	}
}

type TransactionService interface {
	Withdraw(context.Context, *model.Transaction) error
	Deposit(context.Context, *model.Transaction) error
	Transfer(context.Context, *model.Transaction) error
	// PollScheduledTransaction(ctx context.Context, time time.Time) (err error)
	// SaveScheduledTransaction(ctx context.Context, transaction *model.ScheduledTransaction) (err error)
	// ConsumeScheduledTransaction(ctx context.Context) (err error)
	GetTransferDetail(c context.Context, td *model.TransactionDetail) (err error)
	GetAllTransactionByAccountNo(c context.Context, request model.TransactionHistoryRequest) (*[]model.TransactionHistoryResponse, error)
}
