package service

import (
	"context"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (a *transactionService) createTransaction(ctx context.Context, tr *model.Transaction) (err error) {
	if err = a.transactionRepo.CreateTransaction(ctx, tr); err != nil {
		return err
	}
	return nil
}
