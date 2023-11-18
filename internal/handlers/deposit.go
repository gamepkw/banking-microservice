package handler

import (
	"fmt"
	"net/http"
	"time"

	model "github.com/gamepkw/transactions-banking-microservice/internal/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (a *TransactionHandler) Deposit(c echo.Context) (err error) {

	time.Sleep(2 * time.Second)
	var transaction model.Transaction
	transaction.SubmittedAt = time.Now()
	validate := validator.New()

	if err = c.Bind(&transaction); err != nil {

		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := validate.Struct(&transaction); err != nil {
		fmt.Println(err)
		return err
	}

	if transaction.Type != "deposit" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	ctx := c.Request().Context()

	if err = a.transactionService.Deposit(ctx, &transaction); err != nil {
		return c.JSON(getStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, TransactionResponse{Message: "Deposit successfully", Body: &transaction})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case model.ErrInternalServerError:
		return http.StatusInternalServerError
	case model.ErrNotFound:
		return http.StatusNotFound
	case model.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
