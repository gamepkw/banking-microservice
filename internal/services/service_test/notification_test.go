package service

import (
	"context"
	"fmt"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	producer "github.com/gamepkw/transactions-banking-microservice/pkg/kafka/producer"
	"github.com/spf13/viper"
)

func (a *mockTransactionService) migrateTransactionHistoryResponse(ctx context.Context) (err error) {
	return nil
}

func (a *mockTransactionService) addTransactionNotiToQueue(ctx context.Context, tr model.Transaction, remainingBalance float64) {
	topic := "sms_transaction"
	brokerAddress := viper.GetString("kafka.broker_address")
	if tr.Type == "withdraw" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)
		producer.RunKafkaProducer(brokerAddress, topic, message)
	} else if tr.Type == "deposit" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)

		producer.RunKafkaProducer(brokerAddress, topic, message)
	} else if tr.Type == "transfer" {
		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f|%s",
			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance, tr.Receiver.AccountNo)

		producer.RunKafkaProducer(brokerAddress, topic, message)
	}
}
