package middlewares

import (
	//"fmt"
	//"bankstore/api"
	"errors"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	//Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	return
}

/*
	func GenerateJWT(email string, username string) (tokenString string, err error) {
		expirationTime := time.Now().Add(1 * time.Hour)
		claims:= &JWTClaim{
			Email: email,
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString(jwtKey)
		return
	}
*/
func ValidateToken(signedToken string) (user string, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	user = claims.Username

	return

}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user string

		tokenString, err := context.Cookie("Authorization")
		if err != nil {
			//context.JSON(401, gin.H{"error!!!!!!": err.Error()})
			context.Redirect(http.StatusFound, "/login")

			context.Abort()
			return
		}

		if tokenString == "" {
			//context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Redirect(http.StatusFound, "/login")
			context.Abort()
			return
		}

		user, err = ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Redirect(http.StatusFound, "/login")
			context.Abort()
			return
		}
		context.Set("user", user)

		context.Next()
	}
}
