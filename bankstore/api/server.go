package api

import (
	"github.com/gin-gonic/gin"
	db "bankstore/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// TODO: add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts/:limit/:offset", server.getListAccounts)
	router.PUT("/account/:id/:balance", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
