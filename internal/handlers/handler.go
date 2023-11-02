package handler

import (
	"github.com/gamepkw/transactions-banking-microservice/internal/middleware"
	"github.com/labstack/echo/v4"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	transactionService "github.com/gamepkw/transactions-banking-microservice/internal/services"

	"github.com/go-redis/redis"

	validator "github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
	transactionService transactionService.TransactionService
	redis              *redis.Client
	validator          *validator.Validate
}

type TransactionResponse struct {
	Message string             `json:"message"`
	Body    *model.Transaction `json:"body,omitempty"`
}

type TransactionDetailResponse struct {
	Message string                   `json:"message"`
	Body    *model.TransactionDetail `json:"body,omitempty"`
}

type TransactionHistoryResponseResponse struct {
	Message string                              `json:"message"`
	Body    *[]model.TransactionHistoryResponse `json:"body,omitempty"`
}

func NewTransactionHandler(e *echo.Echo, us transactionService.TransactionService, redis *redis.Client) {
	handler := &TransactionHandler{
		transactionService: us,
		redis:              redis,
	}

	middL := middleware.InitMiddleware()

	transactionapiGroup := e.Group("/transaction", middL.RateLimitMiddlewareForTransaction)

	// e.GET("/transactions", handler.GetAllTransaction)
	transactionapiGroup.POST("/deposit", handler.Deposit)
	transactionapiGroup.POST("/withdraw", handler.Withdraw)
	transactionapiGroup.POST("/transfer", handler.Transfer)
	transactionapiGroup.POST("/get-transaction-detail", handler.GetTransferDetail)
	transactionapiGroup.POST("/schedule", handler.SetScheduledTransaction)
	transactionapiGroup.POST("/get-all-transaction", handler.GetAllTransactionByAccountNo)
}
