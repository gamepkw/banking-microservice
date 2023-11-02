package repository

import (
	"context"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (m *transactionRepository) CreateScheduledTransaction(ctx context.Context, tr *model.ScheduledTransaction) (err error) {
	query := `
			INSERT INTO banking.scheduled_transactions 
			SET amount=?, type=?, account=?, receiver=?, submitted_at=?, created_at=? , scheduled_execution_at =?, updated_at=?
	`

	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	tr.CreatedAt = time.Now()
	tr.UpdatedAt = time.Now()

	res, err := stmt.ExecContext(ctx, tr.Amount, tr.Type, tr.Account.AccountNo, tr.Receiver.AccountNo, tr.SubmittedAt, tr.CreatedAt, tr.ScheduledExecutionAt, tr.UpdatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	tr.Id = lastID

	return nil
}
