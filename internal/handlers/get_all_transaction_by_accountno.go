package handler

import (
	"net/http"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *TransactionHandler) GetAllTransactionByAccountNo(c echo.Context) error {
	time.Sleep(2 * time.Second)
	var transaction model.TransactionHistoryRequest

	// requestBody := utils.UnmarshalRequestBody(c.Request())

	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	res, err := a.transactionService.GetAllTransactionByAccountNo(ctx, transaction)
	if err != nil {
		// logger.Error(fmt.Sprintf("%s %s \n %s", transferRequest, err.Error(), requestBody), c.Request())
		return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, TransactionHistoryResponseResponse{Message: "Get all transaction history success", Body: res})
}
