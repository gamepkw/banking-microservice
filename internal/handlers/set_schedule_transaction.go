package handler

import (
	"net/http"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *TransactionHandler) SetScheduledTransaction(c echo.Context) error {
	var transaction model.ScheduledTransaction

	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if transaction.Type != "transfer" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if transaction.Account.AccountNo == transaction.Receiver.AccountNo {
		return echo.NewHTTPError(http.StatusBadRequest, "Can not transfer to the same account")
	}

	if transaction.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Transfer amount must be positive")
	}

	// ctx := c.Request().Context()
	transaction.SubmittedAt = time.Now()

	// if err := a.transactionService.SaveScheduledTransaction(ctx, &transaction); err != nil {
	// 	return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	// }

	return c.JSON(http.StatusCreated, TransactionResponse{Message: "Set scheduled transfer successfully", Body: nil})
}
