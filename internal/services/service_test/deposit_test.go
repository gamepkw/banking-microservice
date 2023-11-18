package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	repo "github.com/gamepkw/transactions-banking-microservice/internal/repositories/repositories_test"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var (
	now = time.Now()
)

func TestDeposit(t *testing.T) {
	mockTransactionRepo := repo.NewMockTransactionRepository(NewMockDB(t), repo.NewMockRedisClient())
	mockService := NewMockTransactionService(mockTransactionRepo, 100*time.Millisecond, NewMockRedisClient())

	tr := &model.Transaction{
		Amount: 1000,
		Type:   "withdraw",
		Fee:    0,
		Total:  0,
		Account: accountModel.Account{
			AccountNo: "0955054682",
		},
	}
	ctx := context.TODO()

	err := mockService.Deposit(ctx, tr)
	t.Logf("got: %v", tr)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// You can add additional assertions or validations as needed
}

func (m *mockTransactionService) Deposit(c context.Context, tr *model.Transaction) error {
	ctx := context.TODO()
	mockRedis := NewMockRedisClient()

	acc, err := m.restGetAccountByAccountNo(ctx, tr.Account.AccountNo)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error get account data: %s", tr.Account.AccountNo))
	}

	var limitAmount string

	cacheKey := "min_deposit_amount"
	limitAmount, err = mockRedis.Get(cacheKey)
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("redis key %s is missing\n", cacheKey)
			limitAmount, _ = m.mockTransactionRepo.GetTransactionConfig(ctx, cacheKey)
		}
		if err := mockRedis.Set(cacheKey, limitAmount, 100*time.Second); err != nil {
			return err
		}
	}

	fmt.Println(limitAmount)

	floatLimitAmount, _ := strconv.ParseFloat(limitAmount, 64)

	if floatLimitAmount > tr.Amount {
		return model.ErrMinimumDeposit
	}
	tr.Total = tr.Amount
	acc.Balance += tr.Total

	if err = m.restUpdateAccount(ctx, *acc); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error update account: %s", acc.AccountNo))
	}

	if err = m.createTransaction(ctx, tr); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error create transaction"))
	}

	tr.Account = *acc

	// go m.addTransactionNotiToQueue(ctx, *tr, acc.Balance)

	return nil

}

// func (m *mockTransactionService) restGetAccountByAccountNo(ctx context.Context, accountNo string) (*accountModel.Account, error) {
// 	var accountJsonString string
// 	if accountNo == "0955054682" {
// 		accountJsonString = `{
// 			"account_no": "0955054682",
// 			"uuid": "2cada3e9174a5fdd8cf26e82ed1829a3f9ee15d8f2fe3d21ee29d12b61997af9",
// 			"balance": 200000,
// 			"bank": "KTB",
// 			"status": "active",
// 			"created_at": "2023-09-21T10:23:59.129859+07:00",
// 			"updated_at": "2023-11-07T11:36:34.52307+07:00"
// 		}`
// 	} else if accountNo == "0163442066" {
// 		accountJsonString = `{
// 			"account_no": "0163442066",
// 			"uuid": "18c11f6b761730a4778c914118f86d7802cf1ffd9f9b5107cfe477d56ce7b176",
// 			"balance": 200000,
// 			"bank": "GSB",
// 			"status": "active",
// 			"created_at": "2023-10-15T18:34:32.38655+07:00",
// 			"updated_at": "2023-11-06T23:50:42.20126+07:00"
// 		}`
// 	}

// 	bytes := []byte(accountJsonString)
// 	var account accountModel.Account
// 	if err := json.Unmarshal(bytes, &account); err != nil {
// 		return nil, errors.Wrap(err, "error unmarshal response body")
// 	}

// 	return &account, nil
// }

func (m *mockTransactionService) restGetAccountByAccountNo(ctx context.Context, accountNo string) (*accountModel.Account, error) {
	args := m.Called(ctx, accountNo)
	return args.Get(0).(*accountModel.Account), args.Error(1)
}

// func (m *mockTransactionService) restUpdateAccount(ctx context.Context, account accountModel.Account) error {
// 	// *account.UpdatedAt = now
// 	return nil
// }

func (m *mockTransactionService) restUpdateAccount(ctx context.Context, account accountModel.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}
