package service

import (
	"context"
	"time"

	// producer "github.com/atm5_microservices/kafka/producer"

	// "github.com/IBM/sarama"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	transactionRepo "github.com/gamepkw/transactions-banking-microservice/internal/repositories"

	"github.com/go-redis/redis"
)

type transactionService struct {
	transactionRepo transactionRepo.TransactionRepository
	contextTimeout  time.Duration
	redis           *redis.Client
	// kafkaClient     sarama.Client
}

func NewTransactionService(
	tr transactionRepo.TransactionRepository,
	timeout time.Duration,
	redis *redis.Client,
	// kafka sarama.Client
) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		contextTimeout:  timeout,
		redis:           redis,
		// kafkaClient:     kafka,
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
