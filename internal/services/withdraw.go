package service

import (
	"context"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (a *transactionService) Withdraw(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	acc, err := a.restGetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return err
	}

	if acc.Balance < tr.Amount {
		return model.ErrInsufficientBalance
	}
	tr.Total = tr.Amount
	acc.Balance -= tr.Total

	if err = a.restUpdateAccount(ctx, *acc); err != nil {
		return err
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return err
	}

	tr.Account = *acc

	// go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}
