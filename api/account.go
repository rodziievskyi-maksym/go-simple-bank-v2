package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"

	_ "github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var request createAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
	}

	account, err := s.store.CreateAccount(c, arg)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	Limit int32 `form:"limit" binding:"required,min=5,max=10"`
	Page  int32 `form:"page" binding:"required,min=1"`
}

func (r *listAccountsRequest) offset() int32 {
	return (r.Page - 1) * r.Limit
}

func (s *Server) listAccounts(c *gin.Context) {
	var request listAccountsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  request.Limit,
		Offset: request.offset(),
	}

	accounts, err := s.store.ListAccounts(c, arg)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

type getAccountRequest struct {
	ID int64 `uri:"id" json:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(c *gin.Context) {
	var request getAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(c, request.ID)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, account)
}
