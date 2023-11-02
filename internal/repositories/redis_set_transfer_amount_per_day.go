package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/go-redis/redis"
)

func (m *transactionRepository) SetTransferAmountPerDayInRedis(ctx context.Context, tr *model.Transaction) error {
	cacheKey := fmt.Sprintf("daily_transaction_%s", tr.Account.AccountNo)
	cachedAmount, err := m.redis.Get(cacheKey).Result()
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	expirationTime := nextMidnight
	durationUntilExpiration := expirationTime.Sub(now)

	if err == redis.Nil {
		err = m.redis.Set(cacheKey, tr.Amount, durationUntilExpiration).Err()
		if err != nil {
			return nil
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("error parsing user from cache: %v", err)
	} else {
		amount := cachedAmount
		floatamount, _ := strconv.ParseFloat(amount, 64)
		floatamount += tr.Amount
		err = m.redis.Set(cacheKey, floatamount, 0).Err()
		if err != nil {
			return fmt.Errorf("error parsing user from cache: %v", err)
		}
		return nil
	}

}
