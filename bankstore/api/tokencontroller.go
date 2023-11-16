package api

import (
	"bankstore/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type TokenRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Username   string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	RePassword string `form:"repassword" json:"repassword" binding:"required"`
}


func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	//fmt.Println("providedPassword:", providedPassword, "----", "user.Password:", user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) GenerateToken(context *gin.Context) {
	var form TokenRequest
	var user User

	context.SetCookie("Authorization", "", -1, "/", context.Request.URL.Hostname(), false, true) //4 Logout

	if err := context.ShouldBind(&form); err != nil {
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Bankstore system",
			//"message": "Error - bad request",
		})
		context.Abort()
		return
	}

	if form.Password == "" {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": "password is empty!"})
		//fmt.Println("Password is empty!!!")
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": "Error - password is empty",
		})
		context.Abort()
		return
	}

	// check if users exists and password is correct
	dbUser, err := server.store.GetUserByName(context, form.Username)
	if err != nil {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		//context.Redirect(http.StatusMovedPermanently, "/login")
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": "Error - username not found",
		})
		context.Abort()
		return
	}

	user = User{
		Username:   dbUser.Username,
		Password:   dbUser.Password,
		RePassword: dbUser.Repassword,
	}

	fmt.Println("user", user.Password, "form.Password", form.Password)
	credentialError := user.CheckPassword(form.Password)
	if credentialError != nil {
		//context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		//context.Redirect(http.StatusMovedPermanently, "/login")
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": "Error - invalid credentials",
		})
		context.Abort()
		return
	}

	//tokenString, err:= auth.GenerateJWT(user.Email, user.Username)
	tokenString, err := middlewares.GenerateJWT(user.Username)
	if err != nil {
		//context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//context.Redirect(http.StatusMovedPermanently, "/login")
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":   "Bankstore system",
			"message": "Error - generate JWT error",
		})
		context.Abort()
		return
	}
	context.SetCookie("Authorization", tokenString, 3600, "/", context.Request.URL.Hostname(), false, true)
	context.Redirect(http.StatusFound, "/index") // "/secured/index"

}
