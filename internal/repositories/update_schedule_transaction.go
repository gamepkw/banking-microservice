package repository

import (
	"context"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (m *transactionRepository) UpdateScheduledTransaction(ctx context.Context, tr model.ScheduledTransaction) (err error) {
	query := `UPDATE banking.scheduled_transactions set status=?, updated_at=? WHERE id = ?`

	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	tr.UpdatedAt = time.Now()

	_, err = stmt.ExecContext(ctx, tr.Status, tr.UpdatedAt, tr.Id)
	if err != nil {
		return
	}
	return
}
