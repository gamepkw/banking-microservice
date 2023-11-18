package service

import (
	"context"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (a *mockTransactionService) createTransaction(ctx context.Context, tr *model.Transaction) (err error) {
	return nil
}
