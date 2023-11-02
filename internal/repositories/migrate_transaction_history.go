package repository

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (m *transactionRepository) MigrateTransactionHistoryResponse(ctx context.Context) (err error) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	log.Printf("Starting migration for date: %s", currentDate)

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	log.Printf("Counting rows to be migrated...")

	select_query := `
					SELECT COUNT(*) FROM banking.transactions 
					WHERE created_at < ?;
	`

	stmt, err := m.conn.PrepareContext(ctx, select_query)
	if err != nil {
		return err
	}

	var countRows int64
	err = stmt.QueryRowContext(ctx, currentDate).Scan(&countRows)
	if err != nil {
		return err
	}

	migrate_query := `
					INSERT INTO banking.transactions_history
					SELECT * FROM banking.transactions
					WHERE created_at < ?;
	`

	stmt, err = m.conn.PrepareContext(ctx, migrate_query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, currentDate)
	if err != nil {
		return err
	}

	insertRowAffect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countRows != insertRowAffect {
		tx.Rollback()
		return fmt.Errorf("number of rows affected does not match")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	delete_query := `
					DELETE FROM banking.transactions
					WHERE created_at < ?;
	`

	stmt, err = m.conn.PrepareContext(ctx, delete_query)
	if err != nil {
		return
	}

	res, err = stmt.ExecContext(ctx, currentDate)
	if err != nil {
		return err
	}

	deleteRowAffect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deleteRowAffect != insertRowAffect {
		tx.Rollback()
		return fmt.Errorf("number of rows affected in delete and migrate operations do not match")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Migration completed successfully")

	return nil
}
