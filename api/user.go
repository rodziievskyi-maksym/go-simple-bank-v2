package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"
	"github.com/rodziievskyi-maksym/go-simple-bank-v2/util"
)

type createUserRequest struct {
	//alphanum validation means it should contain ascii values only
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (s *Server) createUser(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	}

	user, err := s.store.CreateUser(c, arg)
	if err != nil {
		statusCode, response := handleDatabaseError(err)
		c.JSON(statusCode, response)
		return
	}

	response := createUserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
