package handler

import (
	"net/http"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *TransactionHandler) Transfer(c echo.Context) error {
	time.Sleep(2 * time.Second)
	// logger.Info(fmt.Sprintf("%s: start...", transferRequest), c.Request())
	var transaction model.Transaction

	// requestBody := utils.UnmarshalRequestBody(c.Request())

	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if transaction.Type != "transfer" {
		// logger.Error(fmt.Sprintf("%s: Invalid request \n %s", transferRequest, requestBody), c.Request())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if transaction.Account.AccountNo == transaction.Receiver.AccountNo {
		// logger.Error(fmt.Sprintf("%s: Can not transfer to the same account \n %s", transferRequest, requestBody), c.Request())
		return echo.NewHTTPError(http.StatusBadRequest, "Can not transfer to the same account")
	}

	if transaction.Amount <= 0 {
		// logger.Error(fmt.Sprintf("%s: Transfer amount must be positive \n %s", transferRequest, requestBody), c.Request())
		return echo.NewHTTPError(http.StatusBadRequest, "Transfer amount must be positive")
	}

	ctx := c.Request().Context()
	transaction.SubmittedAt = time.Now()

	if err := a.transactionService.Transfer(ctx, &transaction); err != nil {
		// logger.Error(fmt.Sprintf("%s %s \n %s", transferRequest, err.Error(), requestBody), c.Request())
		return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	// logger.Info(fmt.Sprintf("%s: stop...", transferRequest), c.Request())
	return c.JSON(http.StatusCreated, TransactionResponse{Message: "Transfer successfully", Body: &transaction})
}
