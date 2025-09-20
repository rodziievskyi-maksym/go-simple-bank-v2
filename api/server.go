package api

import (
	"github.com/gin-gonic/gin"
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

	//add routes
	accountsGroup := router.Group("/accounts")
	accountsGroup.POST("/", server.createAccount)
	accountsGroup.GET("/", server.listAccounts)
	accountsGroup.GET("/:id", server.getAccount)

	server.router = router
	return server
}

func (s *Server) ServeHTTP(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
