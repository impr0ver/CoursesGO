package api

import (
	"net/http"
	db "bankstore/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) IndexPageGet(context *gin.Context) {
	loginUser, exists := context.Get("user")
	if !exists {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": "user is not exists"})
		context.Redirect(http.StatusFound, "/login")
		context.Abort()
		return
	}

	arg := db.ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}

	accounts, err := server.store.ListAccounts(context, arg)
	if err != nil {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		//context.Redirect(http.StatusMovedPermanently, "/login")
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": http.StatusInternalServerError,
			
		})
	}

	context.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Bankstore system",
		"user":  loginUser,
		"accounts": accounts,
	})
}
