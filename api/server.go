package api

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	db "github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"
)

type Server struct {
	store  db.Store //db interface
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	_ = router.SetTrustedProxies([]string{"127.0.0.1"})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatalf("failed to register validation: %v", err)
		}
	}

	//add routes
	accountsGroup := router.Group("/accounts")
	accountsGroup.POST("/", server.createAccount)
	accountsGroup.GET("/", server.listAccounts)
	accountsGroup.GET("/:id", server.getAccount)

	transfersGroup := router.Group("/transfers")
	transfersGroup.POST("/", server.createTransfer)

	usersGroup := router.Group("/users")
	usersGroup.POST("/", server.createUser)

	server.router = router
	return server
}

func (s *Server) ServeHTTP(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func handleDatabaseError(err error) (int, gin.H) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		slog.Error("database error",
			"code", pqErr.Code.Name(),
			"message", pqErr.Message,
			"detail", pqErr.Detail,
			"constraint", pqErr.Constraint,
		)

		switch pqErr.Code.Name() {
		case "unique_violation":
			return http.StatusConflict, gin.H{
				"error": "Resource already exists",
				"type":  "conflict",
			}
		case "foreign_key_violation":
			return http.StatusBadRequest, gin.H{
				"error": "Invalid reference",
				"type":  "validation",
			}
		case "check_violation":
			return http.StatusBadRequest, gin.H{
				"error": "Invalid data",
				"type":  "validation",
			}
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, gin.H{
			"error": "Resource not found",
			"type":  "not_found",
		}
	}

	slog.Error("unexpected database error", "error", err)
	return http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
		"type":  "internal",
	}
}
