package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
	"github.com/pkg/errors"
)

func (a *transactionService) restGetAccountByAccountNo(ctx context.Context, accountNo string) (*accountModel.Account, error) {
	httpClient := &http.Client{}
	getAccountUrl := fmt.Sprintf("http://localhost:8070/accounts/%s", accountNo)
	req, err := http.NewRequest("GET", getAccountUrl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error make http request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error get http response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve account details, status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error read response body")
	}

	var account accountModel.Account
	if err := json.Unmarshal(responseBody, &account); err != nil {
		return nil, errors.Wrap(err, "error unmarshal response body")
	}

	return &account, nil
}

func (a *transactionService) restUpdateAccount(ctx context.Context, account accountModel.Account) error {
	httpClient := &http.Client{}
	updateAccountURL := "http://localhost:8070/accounts/update"

	requestBody, err := json.Marshal(account)
	if err != nil {
		return errors.Wrap(err, "error marshal request body")
	}

	req, err := http.NewRequest("PUT", updateAccountURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.Wrap(err, "error make http request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error get http response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update account, status code: %d", resp.StatusCode)
	}

	var responseBodyBytes []byte
	responseBodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error read response body")
	}

	if len(responseBodyBytes) == 0 {
		return fmt.Errorf("empty response body")
	}

	if err := json.Unmarshal(responseBodyBytes, &account); err != nil {
		return errors.Wrap(err, "error unmarshal response body")
	}

	return nil
}

func (a *transactionService) restGetDailyRemainingAmount(ctx context.Context, accountNo string) (float64, error) {
	httpClient := &http.Client{}

	getDailyRemainingAmountUrl := fmt.Sprintf("http://localhost:8070/accounts/get-daily-remaining-amount/%s", accountNo)

	httpRequest, err := http.NewRequest("GET", getDailyRemainingAmountUrl, nil)

	if err != nil {
		return 0, err
	}
	// bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWwiOiIyY2FkYTNlOTE3NGE1ZmRkOGNmMjZlODJlZDE4MjlhM2Y5ZWUxNWQ4ZjJmZTNkMjFlZTI5ZDEyYjYxOTk3YWY5IiwiZXhwIjoxNjk5Mjg5NjQyfQ.-n2iF1RZOwRHv2lEBDsWL7lop3Z__9MHfj7vsQ0g9s4"
	// httpRequest.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to retrieve account details, status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var dailyRemainingAmount float64
	if err := json.Unmarshal(responseBody, &dailyRemainingAmount); err != nil {
		return 0, err
	}

	return dailyRemainingAmount, nil
}
