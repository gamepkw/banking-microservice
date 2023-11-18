package repository

import (
	"context"
	"fmt"
	"log"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (m *transactionRepository) GetAllTransactionByAccountNo(ctx context.Context, request model.TransactionHistoryRequest) ([]model.TransactionHistoryResponse, error) {

	query := `
			SELECT type, total_amount, receiver, created_at FROM banking.transactions WHERE account = ? 
	`

	if request.Filter.Type != "" {
		query += fmt.Sprintf("AND type = '%s' ", request.Filter.Type)
	}
	if request.Filter.Year != "" && request.Filter.Month != "" {
		query += fmt.Sprintf("AND DATE_FORMAT(created_at, '%%Y-%%m') = '%s-%s' ", request.Filter.Year, request.Filter.Month)
	}
	if request.Filter.Year != "" && request.Filter.Month == "" {
		query += fmt.Sprintf("AND DATE_FORMAT(created_at, '%%Y') = '%s' ", request.Filter.Year)
	}
	if request.Filter.Year == "" && request.Filter.Month != "" {
		query += fmt.Sprintf("AND DATE_FORMAT(created_at, '%%Y-%%m') = '2023-%s' ", request.Filter.Month)
	}

	query += "ORDER BY created_at DESC"

	rows, err := m.conn.Query(query, request.AccountNo)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	transactions := []model.TransactionHistoryResponse{}
	for rows.Next() {
		transaction := model.TransactionHistoryResponse{}
		err := rows.Scan(
			&transaction.Type,
			&transaction.Total,
			&transaction.Receiver,
			&transaction.CreatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		transactions = append(transactions, transaction)

	}

	return transactions, nil
}
