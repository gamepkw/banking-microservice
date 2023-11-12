package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/pkg/errors"
)

func (a *mockTransactionService) Withdraw(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	acc, err := a.restGetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error get account data: %s", tr.Account.AccountNo))
	}

	if acc.Balance < tr.Amount {
		return model.ErrInsufficientBalance
	}
	tr.Total = tr.Amount
	acc.Balance -= tr.Total

	if err = a.restUpdateAccount(ctx, *acc); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error update account: %s", acc.AccountNo))
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error create transaction: "))
	}

	tr.Account = *acc

	go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}
