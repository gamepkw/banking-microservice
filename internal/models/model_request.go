package model

import (
	"time"

	accountModel "github.com/gamepkw/accounts-banking-microservice/models"
)

type Transaction struct {
	Id          int64                `json:"id"`
	Amount      float64              `json:"amount"`
	Type        string               `json:"type"`
	Fee         float64              `json:"fee"`
	Total       float64              `json:"total"`
	SubmittedAt time.Time            `json:"submitted_at"`
	CreatedAt   time.Time            `json:"created_at"`
	Account     accountModel.Account `json:"account"`
	Receiver    accountModel.Account `json:"receiver,omitempty"`
}

type ScheduledTransaction struct {
	Id                   int64                `json:"id"`
	Amount               float64              `json:"amount"`
	Type                 string               `json:"type"`
	Account              accountModel.Account `json:"account"`
	Receiver             accountModel.Account `json:"receiver,omitempty"`
	Status               string               `json:"status"`
	SubmittedAt          time.Time            `json:"submitted_at"`
	CreatedAt            time.Time            `json:"created_at"`
	ScheduledExecutionAt string               `json:"scheduled_execution_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type TransactionDetail struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount" validate:"required,min=0.01"`
	Fee      float64 `json:"fee"`
	Total    float64 `json:"total"`
}

type TransactionHistoryRequest struct {
	AccountNo string                          `json:"account_no"`
	Filter    TransactionHistoryRequestFilter `json:"filter"`
}

type TransactionHistoryRequestFilter struct {
	Type  string `json:"type"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TransactionRequest struct {
	Amount   float64              `json:"amount"`
	Type     string               `json:"type"`
	Account  accountModel.Account `json:"account"`
	Receiver accountModel.Account `json:"receiver,omitempty"`
}
