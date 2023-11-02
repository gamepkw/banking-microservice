package service

import (
	"context"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/pkg/errors"
)

func (a *transactionService) GetAllTransactionByAccountNo(c context.Context, request model.TransactionHistoryRequest) (*[]model.TransactionHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.transactionRepo.GetAllTransactionByAccountNo(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all accounts")
	}

	return &res, nil
}
