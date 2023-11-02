package repository

import (
	"context"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
)

func (m *transactionRepository) CreateTransaction(ctx context.Context, tr *model.Transaction) (err error) {
	query := `
			INSERT INTO banking.transactions 
			SET amount=?, type=?, fee=?, total_amount=?, submitted_at=?, created_at=? , account=?, receiver=?
	`

	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	tr.CreatedAt = time.Now()

	res, err := stmt.ExecContext(ctx, tr.Amount, tr.Type, tr.Fee, tr.Total, tr.SubmittedAt, tr.CreatedAt, tr.Account.AccountNo, tr.Receiver.AccountNo)
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
