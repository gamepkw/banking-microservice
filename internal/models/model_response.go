package model

import "time"

type TransactionHistoryResponse struct {
	Type      string    `json:"type"`
	Total     float64   `json:"total"`
	Receiver  string    `json:"receiver,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
