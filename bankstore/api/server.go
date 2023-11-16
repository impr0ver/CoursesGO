package api

import (
	db "bankstore/db/sqlc"
	"bankstore/middlewares"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	
	router.LoadHTMLGlob("./templates/*.tmpl")
	router.Static("./assets", "./assets")

	router.GET("/login", server.GenerateToken) //old token
	router.POST("/login", server.GenerateToken)

	router.GET("/register", server.RegisterUser)
	router.POST("/register", server.RegisterUser)

	secured := router.Group("").Use(middlewares.Auth()) // /secured middleware 
	{
		secured.GET("/index", server.IndexPageGet)
		secured.GET("", server.IndexPageGet)
	}


	// API:
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
