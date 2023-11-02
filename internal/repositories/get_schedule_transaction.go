package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (m *transactionRepository) GetScheduledTransaction(ctx context.Context, fetchTime time.Time) (transactions []model.ScheduledTransaction, err error) {

	query := `
			SELECT * FROM scheduled_transactions WHERE (scheduled_execution_at = ? AND status = 'unprocessed')
	`

	rows, err := m.conn.Query(query, fetchTime)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		transaction := model.ScheduledTransaction{}
		err := rows.Scan(
			&transaction.Id,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Account.AccountNo,
			&transaction.Receiver.AccountNo,
			&transaction.Status,
			&transaction.SubmittedAt,
			&transaction.CreatedAt,
			&transaction.ScheduledExecutionAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		transactions = append(transactions, transaction)

		transaction.Status = "processing"
		m.UpdateScheduledTransaction(ctx, transaction)

	}

	return
}
