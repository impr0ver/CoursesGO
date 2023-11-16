package api

import (
	//"jwt-authentication-golang/database"
	//"jwt-authentication-golang/models"
	db "bankstore/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) RegisterUser(context *gin.Context) {
	var user User
	//var form TokenRequest

	if err := context.ShouldBind(&user); err != nil {
		//context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "Hashes system",
			//"message": "Error bind form",
		})
		context.Abort()
		return
	}

	if user.Username == "" {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": "password is empty!"})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Hashes system",
			"message": "Error - username is empty",
		})
		context.Abort()
		return
	}

	if user.Password == "" {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": "password is empty!"})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Hashes system",
			"message": "Error - password is empty",
		})
		context.Abort()
		return
	}

	if user.Password != user.RePassword {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": "passwords are not equal!"})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Hashes system",
			"message": "Error - passwords are not equal",
		})
		context.Abort()
		return
	}

	user.RePassword = ""

	if err := user.HashPassword(user.Password); err != nil {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Hashes system",
			"message": "Error - hash password error",
		})
		context.Abort()
		return
	}

	arg := db.CreateUserParams{
		Username:   user.Username,
		Password:   user.Password,
		Repassword: user.RePassword,
	}
	//record := database.Instance.Create(&user)
	_, err := server.store.CreateUser(context, arg)
	if err != nil {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": "Error - create user in db",
		})
		context.Abort()
		return
	}

	//context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
	//context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "username": user.Username})
	context.HTML(http.StatusOK, "register.tmpl", gin.H{
		"title":    "Hashes system",
		"username": user.Username,
	})

}
