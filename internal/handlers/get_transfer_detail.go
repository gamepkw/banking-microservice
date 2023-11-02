package handler

import (
	"net/http"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *TransactionHandler) GetTransferDetail(c echo.Context) error {
	var transaction model.TransactionDetail

	// requestBody := utils.UnmarshalRequestBody(c.Request())

	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if transaction.Sender == transaction.Receiver {
		// logger.Error(fmt.Sprintf("%s: Can not transfer to the same account \n %s", transferRequest, requestBody), c.Request())
		return echo.NewHTTPError(http.StatusBadRequest, "Can not transfer to the same account")
	}

	if transaction.Amount <= 0 {
		// logger.Error(fmt.Sprintf("%s: Transfer amount must be positive \n %s", transferRequest, requestBody), c.Request())
		return echo.NewHTTPError(http.StatusBadRequest, "Transfer amount must be positive")
	}

	ctx := c.Request().Context()

	if err := a.transactionService.GetTransferDetail(ctx, &transaction); err != nil {
		// logger.Error(fmt.Sprintf("%s %s \n %s", transferRequest, err.Error(), requestBody), c.Request())
		return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	// logger.Info(fmt.Sprintf("%s: stop...", transferRequest), c.Request())
	return c.JSON(http.StatusOK, TransactionDetailResponse{Message: "Get transaction detail success", Body: &transaction})
}
