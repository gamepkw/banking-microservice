package service

import (
	"context"
)

func (a *transactionService) migrateTransactionHistoryResponse(ctx context.Context) (err error) {
	if err = a.transactionRepo.MigrateTransactionHistoryResponse(ctx); err != nil {
		return err
	}
	return nil
}

// func (a *transactionService) addTransactionNotiToQueue(ctx context.Context, tr model.Transaction, remainingBalance float64) {
// 	topic := "sms_transaction"
// 	brokerAddress := viper.GetString("kafka.broker_address")
// 	if tr.Type == "withdraw" {
// 		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
// 			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)
// 		producer.RunKafkaProducer(brokerAddress, topic, message)
// 	} else if tr.Type == "deposit" {
// 		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f",
// 			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance)

// 		producer.RunKafkaProducer(brokerAddress, topic, message)
// 	} else if tr.Type == "transfer" {
// 		message := fmt.Sprintf("%s|%.2f|%s|%s|%.2f|%s",
// 			tr.Type, tr.Amount, tr.Account.AccountNo, tr.CreatedAt.Format("2006-01-02 15:04:05"), remainingBalance, tr.Receiver.AccountNo)

// 		producer.RunKafkaProducer(brokerAddress, topic, message)
// 	}
// }
