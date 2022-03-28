package middleware

import (
	"fmt"
	"net/http"
	"p4_web/constant"
	db "p4_web/database"
	"p4_web/model"
	"p4_web/tools/env"
	"p4_web/tools/exception"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := strings.Split(c.Request.Header["Authorization"][0], " ")
		if s[0] == "Bearer" {
			token, err := jwt.ParseWithClaims(s[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(env.Get("SECRET_KEY")), nil
			})
			if err != nil {
				panic(exception.ApiException{
					Code:    []int{http.StatusUnauthorized},
					Message: fmt.Sprintf("invalid access token: %v", err.Error()),
				})
			}
			claims := token.Claims.(*jwt.StandardClaims)

			var user model.User
			db.DB.Where("id = ?", claims.Issuer).First(&user)
			c.Set(constant.AUTH, user)
			c.Next()
		} else {
			panic(exception.ApiException{
				Code:    []int{http.StatusUnauthorized},
				Message: "invalid access token",
			})
		}
	}
}
