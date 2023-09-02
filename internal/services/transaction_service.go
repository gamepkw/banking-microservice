package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	producer "github.com/atm5_microservices/kafka/producer"

	"github.com/IBM/sarama"
	accountModel "github.com/gamepkw/accounts-banking-microservice/internal/models"
	model "github.com/gamepkw/transactions-banking-microservice/internal/models"

	// accountModel "github.com/atm5_microservices/accounts_service/internal/models"

	// "github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

type transactionService struct {
	transactionRepo model.TransactionRepository
	accountService  accountModel.AccountService
	contextTimeout  time.Duration
	redis           *redis.Client
	// kafkaClient     sarama.Client
}

// NewTransactionService will create new an transactionService object representation of model.TransactionService interface
func NewTransactionService(tr model.TransactionRepository,
	au accountModel.AccountService,
	timeout time.Duration,
	redis *redis.Client,
	kafka sarama.Client) model.TransactionService {
	return &transactionService{
		transactionRepo: tr,
		accountService:  au,
		contextTimeout:  timeout,
		redis:           redis,
		// kafkaClient:     kafka,
	}
}

// func (a *transactionService) GetTransactionByID(c context.Context, id int64) (res model.Transaction, err error) {
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()

// 	res, err = a.transactionRepo.GetTransactionByTID(ctx, id)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

func (a *transactionService) Withdraw(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	acc, err := a.accountService.GetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return err
	}

	if acc.Balance < tr.Amount {
		return model.ErrInsufficientBalance
	}

	acc.Balance -= tr.Amount

	if err = a.accountService.UpdateAccount(ctx, acc); err != nil {
		return err
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return err
	}

	tr.Account = *acc

	go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}

func (a *transactionService) Deposit(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	acc, err := a.accountService.GetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return err
	}

	cacheKey := "min_deposit_amount"
	limitAmount, err := a.redis.Get(cacheKey).Result()
	if err != nil {
		return fmt.Errorf("redis key %s is missing", cacheKey)
	}

	floatLimitAmount, _ := strconv.ParseFloat(limitAmount, 64)

	if floatLimitAmount > tr.Amount {
		return model.ErrMinimumDeposit
	}

	acc.Balance += tr.Amount

	if err = a.accountService.UpdateAccount(ctx, acc); err != nil {
		return err
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return err
	}

	tr.Account = *acc

	go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil

}

func (a *transactionService) Transfer(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	isExceedLimit := a.checkTransactionLimit(ctx, *tr)

	if !isExceedLimit {
		return fmt.Errorf("exceed limit per day")
	}

	acc, err := a.accountService.GetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return err
	}

	res_acc, err := a.accountService.GetAccountByAccountNo(ctx, tr.Receiver.AccountNo)
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

	if err = a.accountService.UpdateAccount(ctx, acc); err != nil {
		return err
	}

	if err = a.accountService.UpdateAccount(ctx, res_acc); err != nil {
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

	go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}

func (a *transactionService) checkTransactionLimit(ctx context.Context, tr model.Transaction) bool {
	dailyLimit, err := a.accountService.GetDailyLimit(ctx, tr.Account.AccountNo)
	if err != nil {
		return false
	}

	dailySumTransaction, _ := a.accountService.GetSumDailyTransaction(ctx, tr.Account.AccountNo)

	if tr.Amount+dailySumTransaction > dailyLimit {
		return false
	}

	return true
}

func calculateTransferFee(acc, res_acc *model.Account) float64 {
	if acc.Bank != res_acc.Bank {
		return 10.0
	}
	return 0.0
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

func (a *transactionService) createTransaction(ctx context.Context, tr *model.Transaction) (err error) {
	if err = a.transactionRepo.CreateTransaction(ctx, tr); err != nil {
		return err
	}
	return nil
}

func (a *transactionService) migrateTransactionHistory(ctx context.Context) (err error) {
	if err = a.transactionRepo.MigrateTransactionHistory(ctx); err != nil {
		return err
	}
	return nil
}

func (a *transactionService) addTransactionNotiToQueue(ctx context.Context, tr model.Transaction, remainingBalance float64) {
	topic := "sms_transaction"
	brokerAddress := viper.GetString("kafka.broker_address")
	if tr.Type == "withdraw" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)
		producer.RunKafkaProducer(brokerAddress, topic, message)
	} else if tr.Type == "deposit" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)

		producer.RunKafkaProducer(brokerAddress, topic, message)
	} else if tr.Type == "transfer" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f|%s",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance, tr.Receiver.AccountNo)

		producer.RunKafkaProducer(brokerAddress, topic, message)
	}
}
