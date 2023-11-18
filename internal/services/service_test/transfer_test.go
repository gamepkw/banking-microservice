package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	repo "github.com/gamepkw/transactions-banking-microservice/internal/repositories/repositories_test"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func TestTransfer(t *testing.T) {
	// mockService := new(mockTransactionService)
	mockTransactionRepo := repo.NewMockTransactionRepository(NewMockDB(t), repo.NewMockRedisClient())
	// mockRepo := repo.MockTransactionRepository{}
	mockRedis := NewMockRedisClient()
	mockService := NewMockTransactionService(mockTransactionRepo, 100*time.Microsecond, mockRedis)
	ctx := context.Background()
	tr := &model.Transaction{
		Amount: 1000,
		Type:   "transfer",
		Account: accountModel.Account{
			AccountNo: "0955054682",
		},
		Receiver: accountModel.Account{
			AccountNo: "0163442066",
		},
	}

	sender := &accountModel.Account{
		AccountNo: "0955054682",
		Uuid:      "2cada3e9174a5fdd8cf26e82ed1829a3f9ee15d8f2fe3d21ee29d12b61997af9",
		Balance:   200000,
		Bank:      "KTB",
		Status:    "active",
		UpdatedAt: nil,
	}

	receiver := &accountModel.Account{
		AccountNo: "0163442066",
		Uuid:      "18c11f6b761730a4778c914118f86d7802cf1ffd9f9b5107cfe477d56ce7b176",
		Balance:   200000,
		Bank:      "GSB",
		Status:    "active",
		UpdatedAt: nil,
	}

	mockService.On("restGetAccountByAccountNo", mock.AnythingOfType("*context.timerCtx"), tr.Account.AccountNo).Return(sender, nil).Once()
	mockService.On("restGetAccountByAccountNo", mock.AnythingOfType("*context.timerCtx"), tr.Receiver.AccountNo).Return(receiver, nil).Once()

	// mockService.On("restUpdateAccount", mock.AnythingOfType("*context.timerCtx"), sender).
	// 	Return(func(ctx context.Context, sender accountModel.Account) error {
	// 		*sender.UpdatedAt = now
	// 		return nil
	// 	}).Once()
	// mockService.On("restUpdateAccount", mock.AnythingOfType("*context.timerCtx"), receiver).
	// 	Return(func(ctx context.Context, receiver accountModel.Account) error {
	// 		*receiver.UpdatedAt = now
	// 		return nil
	// 	}).Once()

	mockService.On("restUpdateAccount", mock.AnythingOfType("*context.timerCtx"), sender).
		Run(func(args mock.Arguments) {
			argSender := args.Get(1).(*accountModel.Account)
			*argSender.UpdatedAt = now
		}).Return(nil).Once()

	// Mock the restUpdateAccount function for receiver
	mockService.On("restUpdateAccount", mock.AnythingOfType("*context.timerCtx"), receiver).
		Run(func(args mock.Arguments) {
			argReceiver := args.Get(1).(*accountModel.Account)
			*argReceiver.UpdatedAt = now
		}).Return(nil).Once()

	// mockService.On("restUpdateAccount", mock.AnythingOfType("*context.timerCtx"), sender).
	// 	Return(sender, nil).Once()

	err := mockService.Transfer(ctx, tr)
	t.Logf("got: %v", tr)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func (a *mockTransactionService) Transfer(c context.Context, tr *model.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	isExceedLimitPerDay, err := a.checkTransactionLimit(ctx, *tr)
	if err != nil {
		return errors.Wrap(err, "error check transaction limit")
	}

	if !isExceedLimitPerDay {
		return fmt.Errorf("exceed limit per day")
	}

	acc, err := a.restGetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error get sender account data: %s", tr.Account.AccountNo))
	}

	res_acc, err := a.restGetAccountByAccountNo(ctx, tr.Receiver.AccountNo)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error get receiver account data: %s", tr.Receiver.AccountNo))
	}

	if res_acc == nil {
		return model.ErrResipientNotFound
	}

	if res_acc.Status == "inactive" {
		return model.ErrAccDeleted
	}
	if err := a.checkTransferLimit(ctx, acc.AccountNo, tr.Amount); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error check transfer limit account: %s", acc.AccountNo))
	}
	tr.Fee = calculateTransferFee(acc, res_acc)
	tr.Total = tr.Amount + tr.Fee
	if acc.Balance < tr.Total {
		return model.ErrInsufficientBalance
	}
	acc.Balance -= (tr.Total)
	res_acc.Balance += tr.Amount

	if err = a.restUpdateAccount(ctx, *acc); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error update sender account: %s", acc.AccountNo))
	}

	if err = a.restUpdateAccount(ctx, *res_acc); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error update receiver account: %s", res_acc.AccountNo))
	}

	if err = a.createTransaction(ctx, tr); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error create transaction"))
	}

	// if err = a.mockTransactionRepo.SetTransferAmountPerDayInRedis(ctx, tr); err != nil {
	// 	return errors.Wrap(err, fmt.Sprintf("error set transfer amount per day in redis"))
	// }

	tr.Account = *acc
	tr.Receiver = *res_acc

	go a.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil
}

func jsonToStruct(jsonData string) model.Transaction {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(jsonData), &transaction)
	if err != nil {
		fmt.Println("Error:", err)
		return model.Transaction{}
	}
	return transaction
}

func calculateTransferFee(acc, res_acc *accountModel.Account) float64 {
	if acc.Bank != res_acc.Bank {
		return 10.0
	}
	return 0.0
}

func (m *mockTransactionService) checkTransferLimit(ctx context.Context, AccountNo string, amount float64) error {
	mockRedis := NewMockRedisClient()
	cacheKey := fmt.Sprintf("limit_per_transaction: %s", AccountNo)
	limitAmountPerTransaction, err := mockRedis.Get(cacheKey)
	if err == redis.Nil {
		limitAmountPerTransaction = "0"
	}
	limitAmountPerTransactionFloat64, _ := strconv.ParseFloat(limitAmountPerTransaction, 64)

	if limitAmountPerTransactionFloat64 != 0 && amount > limitAmountPerTransactionFloat64 {
		return model.ErrExceedLimitAmountPerTransaction
	}
	return nil
}

func (a *mockTransactionService) checkTransactionLimit(ctx context.Context, tr model.Transaction) (bool, error) {
	dailyRemainingAmount, err := a.restGetDailyRemainingAmount(ctx, tr.Account.AccountNo)
	if err != nil {
		return false, err
	}
	if tr.Amount > dailyRemainingAmount {
		return false, err
	}

	return true, nil
}
func (a *mockTransactionService) restGetDailyRemainingAmount(ctx context.Context, accountNo string) (float64, error) {

	return 100000, nil
}

func (a *mockTransactionService) GetTransferDetail(c context.Context, tr *model.TransactionDetail) (err error) {
	return nil
}

// func TestTransfer(t *testing.T) {

// 	// updatedAt := time.Date(2023, 9, 21, 10, 23, 59, 129859000, time.FixedZone("ICT", 7*60*60))
// 	// updatedAtPointer := &updatedAt
// 	mockService := new(mockTransactionService)
// 	mockTransactionRepo := repo.NewMockTransactionRepository(NewMockDB(t), repo.NewMockRedisClient())
// 	// mockRepo := repo.MockTransactionRepository{}
// 	mockRedis := NewMockRedisClient()
// 	mockService = NewMockTransactionService(mockTransactionRepo, 100*time.Microsecond, mockRedis)
// 	ctx := context.Background()
// 	tr := &model.Transaction{
// 		Amount: 1000,
// 		Type:   "transfer",
// 		Account: accountModel.Account{
// 			AccountNo: "0955054682",
// 		},
// 		Receiver: accountModel.Account{
// 			AccountNo: "0163442066",
// 		},
// 	}

// 	mockService.On("restGetAccountByAccountNo", mock.AnythingOfType("*context.timerCtx"), tr.Account.AccountNo).Return(func(ctx context.Context, accountNo string) (*accountModel.Account, error) {
// 		// Create a new instance with the desired values
// 		account := &accountModel.Account{
// 			AccountNo: "0955054682",
// 			Uuid:      "2cada3e9174a5fdd8cf26e82ed1829a3f9ee15d8f2fe3d21ee29d12b61997af9",
// 			Balance:   200000,
// 			Bank:      "KTB",
// 			Status:    "active",
// 			// UpdatedAt: updatedAtPointer,
// 		}
// 		tr.Account = *account

// 		return account, nil
// 	}).Once()

// 	mockService.On("restGetAccountByAccountNo", mock.AnythingOfType("*context.timerCtx"), tr.Receiver.AccountNo).Return(func(ctx context.Context, accountNo string) (*accountModel.Account, error) {
// 		// Create a new instance with the desired values
// 		account := &accountModel.Account{
// 			AccountNo: "0163442066",
// 			Uuid:      "18c11f6b761730a4778c914118f86d7802cf1ffd9f9b5107cfe477d56ce7b176",
// 			Balance:   200000,
// 			Bank:      "GSB",
// 			Status:    "active",
// 			// UpdatedAt: updatedAtPointer,
// 		}
// 		tr.Receiver = *account

// 		return account, nil
// 	}).Once()

// 	err := mockService.Transfer(ctx, tr)
// 	t.Logf("got: %v", tr)
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}
// }
