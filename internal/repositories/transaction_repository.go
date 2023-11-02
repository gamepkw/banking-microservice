package repository

import (
	"context"
	"database/sql"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"

	"github.com/go-redis/redis"
)

type transactionRepository struct {
	conn  *sql.DB
	redis *redis.Client
}

func NewtransactionRepository(conn *sql.DB, redis *redis.Client) TransactionRepository {
	return &transactionRepository{
		conn:  conn,
		redis: redis,
	}
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tr *model.Transaction) error
	SetTransferAmountPerDayInRedis(ctx context.Context, tr *model.Transaction) error
	MigrateTransactionHistoryResponse(ctx context.Context) (err error)
	CreateScheduledTransaction(ctx context.Context, st *model.ScheduledTransaction) (err error)
	GetScheduledTransaction(ctx context.Context, time time.Time) (transaction []model.ScheduledTransaction, err error)
	UpdateScheduledTransaction(ctx context.Context, tr model.ScheduledTransaction) (err error)
	GetTransactionConfig(ctx context.Context, configName string) (string, error)
	GetAllTransactionByAccountNo(ctx context.Context, request model.TransactionHistoryRequest) ([]model.TransactionHistoryResponse, error)
}
