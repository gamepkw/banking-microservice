package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (a *transactionService) restGetAccountByAccountNo(ctx context.Context, accountNo string) (*accountModel.Account, error) {
	httpClient := &http.Client{}
	getAccountUrl := fmt.Sprintf("http://localhost:8070/users/accounts/%s", accountNo)
	req, err := http.NewRequest("GET", getAccountUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve account details, status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var account accountModel.Account
	if err := json.Unmarshal(responseBody, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *transactionService) restUpdateAccount(ctx context.Context, account accountModel.Account) error {
	httpClient := &http.Client{}
	updateAccountURL := fmt.Sprintf("http://localhost:8070/users/accounts/%s", account.AccountNo) // Fix the URL with "http://"

	requestBody, err := json.Marshal(account)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", updateAccountURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update account, status code: %d", resp.StatusCode)
	}

	var responseBodyBytes []byte
	responseBodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBodyBytes) == 0 {
		return fmt.Errorf("empty response body")
	}

	if err := json.Unmarshal(responseBodyBytes, &account); err != nil {
		return err
	}

	return nil
}

func (a *transactionService) restCheckTransactionLimit(ctx context.Context, req model.TransactionRequest) (bool, error) {
	httpClient := &http.Client{}
	getDailyLimitUrl := fmt.Sprintf("http://localhost:8070/accounts-limit/%s", req.Account.AccountNo)

	httpRequest, err := http.NewRequest("GET", getDailyLimitUrl, nil)
	if err != nil {
		return false, err
	}

	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to retrieve account details, status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dailyLimit float64
	if err := json.Unmarshal(responseBody, &dailyLimit); err != nil {
		return false, err
	}

	getDailySumTransactionUrl := fmt.Sprintf("http://localhost:8070/accounts-daily-limit/%s", req.Account.AccountNo)

	httpRequest, err = http.NewRequest("GET", getDailySumTransactionUrl, nil)
	if err != nil {
		return false, err
	}

	resp, err = httpClient.Do(httpRequest)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to retrieve account details, status code: %d", resp.StatusCode)
	}

	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dailySumTransaction float64
	if err := json.Unmarshal(responseBody, &dailySumTransaction); err != nil {
		return false, err
	}

	if req.Amount+dailySumTransaction > dailyLimit {
		return false, nil
	}

	return true, nil
}
