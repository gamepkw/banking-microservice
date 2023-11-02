package repository

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
)

func (m *transactionRepository) GetTransactionConfig(ctx context.Context, configName string) (string, error) {
	query := `SELECT config_value FROM banking.global_configs WHERE config_name = ?`

	rows, err := m.conn.QueryContext(ctx, query, configName)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	var configValue string

	for rows.Next() {
		err := rows.Scan(&configValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	return configValue, nil
}
