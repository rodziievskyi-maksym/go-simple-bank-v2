package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	server.router = router
	return server
}

func (s *Server) ServeHTTP(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
