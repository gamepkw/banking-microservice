package model

import (
	"context"
	"time"

	accountModel "github.com/gamepkw/accounts-banking-microservice/shared_model"
)

// type Account struct {
// 	AccountNo string     `json:"account_no,omitempty"`
// 	Uuid      string     `json:"uuid,omitempty"`
// 	Name      string     `json:"name,omitempty"`
// 	Email     string     `json:"email,omitempty"`
// 	Tel       string     `json:"tel,omitempty"`
// 	Balance   float64    `json:"balance"`
// 	Bank      string     `json:"bank,omitempty"`
// 	Status    string     `json:"status,omitempty"`
// 	CreatedAt *time.Time `json:"created_at,omitempty"`
// 	UpdatedAt *time.Time `json:"updated_at,omitempty"`
// 	IsClosed  int        `json:"is_closed,omitempty"`
// }

type Transaction struct {
	Id          int64                `json:"id"`
	Amount      float64              `json:"amount"`
	Type        string               `json:"type"`
	Fee         float64              `json:"fee"`
	Total       float64              `json:"total"`
	SubmittedAt time.Time            `json:"submitted_at"`
	CreatedAt   time.Time            `json:"created_at"`
	Account     accountModel.Account `json:"account"`
	Receiver    accountModel.Account `json:"receiver,omitempty"`
}

type ScheduledTransaction struct {
	Id                   int64                `json:"id"`
	Amount               float64              `json:"amount"`
	Type                 string               `json:"type"`
	Account              accountModel.Account `json:"account"`
	Receiver             accountModel.Account `json:"receiver,omitempty"`
	Status               string               `json:"status"`
	SubmittedAt          time.Time            `json:"submitted_at"`
	CreatedAt            time.Time            `json:"created_at"`
	ScheduledExecutionAt string               `json:"scheduled_execution_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type TransactionService interface {
	// GetAllTransaction(ctx context.Context, cursor string, num int64) ([]Transaction, string, error)
	// CreateTransaction(context.Context, *Transaction) error
	Withdraw(context.Context, *Transaction) error
	Deposit(context.Context, *Transaction) error
	Transfer(context.Context, *Transaction) error
	PollScheduledTransaction(ctx context.Context, time time.Time) (err error)
	SaveScheduledTransaction(ctx context.Context, transaction *ScheduledTransaction) (err error)
	ConsumeScheduledTransaction(ctx context.Context) (err error)
}

type TransactionRepository interface {
	// GetAllTransaction(ctx context.Context, cursor string, num int64) (res []Transaction, nextCursor string, err error)
	// GetTransactionByTID(ctx context.Context, tid int64) (Transaction, error)
	CreateTransaction(ctx context.Context, tr *Transaction) error
	SetTransferAmountPerDayInRedis(ctx context.Context, tr *Transaction) error
	MigrateTransactionHistory(ctx context.Context) (err error)
	CreateScheduledTransaction(ctx context.Context, st *ScheduledTransaction) (err error)
	GetScheduledTransaction(ctx context.Context, time time.Time) (transaction []ScheduledTransaction, err error)
	UpdateScheduledTransaction(ctx context.Context, tr ScheduledTransaction) (err error)
}
