package service

import (
	"context"
	"fmt"
	"strconv"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/go-redis/redis"
)

func (a *transactionService) Transfer(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	isExceedLimit, err := a.restCheckTransactionLimit(ctx, model.TransactionRequest{
		Amount:   tr.Amount,
		Type:     tr.Type,
		Account:  tr.Account,
		Receiver: tr.Receiver,
	})
	if err != nil {
		return err
	}

	if !isExceedLimit {
		return fmt.Errorf("exceed limit per day")
	}

	acc, err := a.restGetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return err
	}

	res_acc, err := a.restGetAccountByAccountNo(ctx, tr.Receiver.AccountNo)
	if err != nil {
		return err
	}

	if res_acc == nil {
		return model.ErrResipientNotFound
	}

	if res_acc.Status == "inactive" {
		return model.ErrAccDeleted
	}

	if err := a.checkTransferLimit(ctx, acc.AccountNo, tr.Amount); err != nil {
		return err
	}

	tr.Fee = calculateTransferFee(acc, res_acc)
	tr.Total = tr.Amount + tr.Fee

	if acc.Balance < tr.Total {
		return model.ErrInsufficientBalance
	}

	acc.Balance -= (tr.Total)
	res_acc.Balance += tr.Amount

	if err = a.restUpdateAccount(ctx, *acc); err != nil {
		return err
	}

	if err = a.restUpdateAccount(ctx, *res_acc); err != nil {
		return err
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return err
	}

	if err = a.transactionRepo.SetTransferAmountPerDayInRedis(ctx, tr); err != nil {
		return err
	}

	tr.Account = *acc
	tr.Receiver = *res_acc

	// go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}

func (a *transactionService) checkTransferLimit(ctx context.Context, AccountNo string, amount float64) error {
	cacheKey := fmt.Sprintf("limit_per_transaction: %s", AccountNo)
	limitAmountPerTransaction, err := a.redis.Get(cacheKey).Result()

	if err == redis.Nil {
		limitAmountPerTransaction = "0"
	}
	limitAmountPerTransactionFloat64, _ := strconv.ParseFloat(limitAmountPerTransaction, 64)

	if limitAmountPerTransactionFloat64 != 0 && amount > limitAmountPerTransactionFloat64 {
		return model.ErrExceedLimitAmountPerTransaction
	}
	return nil
}

func calculateTransferFee(acc, res_acc *accountModel.Account) float64 {
	if acc.Bank != res_acc.Bank {
		return 10.0
	}
	return 0.0
}

func (a *transactionService) GetTransferDetail(c context.Context, tr *model.TransactionDetail) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	acc, err := a.restGetAccountByAccountNo(ctx, tr.Sender)
	if err != nil {
		return err
	}

	res_acc, err := a.restGetAccountByAccountNo(ctx, tr.Receiver)
	if err != nil {
		return err
	}

	if res_acc == nil {
		return model.ErrResipientNotFound
	}

	if res_acc.Status == "inactive" {
		return model.ErrAccDeleted
	}

	tr.Fee = calculateTransferFee(acc, res_acc)
	tr.Total = tr.Amount + tr.Fee

	return nil
}
