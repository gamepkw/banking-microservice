package handler

import (
	"net/http"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *TransactionHandler) Withdraw(c echo.Context) (err error) {
	time.Sleep(2 * time.Second)
	var transaction model.Transaction

	ctx := c.Request().Context()

	if err = c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if transaction.Type != "withdraw" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	// ctx := c.Request().Context()

	transaction.SubmittedAt = time.Now()

	if err = a.transactionService.Withdraw(ctx, &transaction); err != nil {
		return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, TransactionResponse{Message: "Withdraw successfully", Body: &transaction})
}
