package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var request createTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !s.validAccount(c, request.FromAccountID, request.Currency) {
		return
	}

	if !s.validAccount(c, request.ToAccountID, request.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	result, err := s.store.TransferTx(c, arg)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] does not match currency [%s] vs [%s]", account.ID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
