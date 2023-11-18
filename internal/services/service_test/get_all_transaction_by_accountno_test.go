package service

import (
	"context"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (a *mockTransactionService) GetAllTransactionByAccountNo(c context.Context, request model.TransactionHistoryRequest) (*[]model.TransactionHistoryResponse, error) {

	return nil, nil
}
