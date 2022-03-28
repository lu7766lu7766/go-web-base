package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"p4_web/config"

	"p4_web/constant"
	db "p4_web/database"
	"p4_web/tools/exception"
	"p4_web/tools/prop"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := strings.Split(c.Request.Header["Authorization"][0], " ")
		if s[0] == "Bearer" {
			token, err := jwt.ParseWithClaims(s[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AuthConfig.SecretKey), nil
			})
			if err != nil {
				panic(exception.ApiException{
					Code:    []int{http.StatusUnauthorized},
					Message: fmt.Sprintf("invalid access token: %v", err.Error()),
				})
			}
			claims := token.Claims.(*jwt.StandardClaims)
			db.DB.Where(fmt.Sprintf("%v = ?", config.AuthConfig.PrimaryKey), claims.Issuer).First(&config.AuthConfig.Model)
			c.Set(constant.AUTH, config.AuthConfig.Model)
			c.Next()
		} else {
			panic(exception.ApiException{
				Code:    []int{http.StatusUnauthorized},
				Message: "invalid access token",
			})
		}
	}
}

func Check(idValue string, password string) {
	query := db.DB
	for index := range config.AuthConfig.Uids {
		fmt.Println(config.AuthConfig.Uids[index])
		query = query.Where(fmt.Sprintf("%v = ?", config.AuthConfig.Uids[index]), idValue)
	}

	err := query.First(&config.AuthConfig.Model).Error

	if err != nil {
		panic(exception.ApiException{
			Code:    []int{1001},
			Message: fmt.Sprintf("user not found: %v", err.Error()),
		})
	}

	if err := bcrypt.CompareHashAndPassword(config.AuthConfig.Model.Password, []byte(password)); err != nil {
		panic(exception.ApiException{
			Code:    []int{1002},
			Message: fmt.Sprintf("password not equal: %v", err.Error()),
		})
	}
}

func Generate() map[string]interface{} {
	user := config.AuthConfig.Model
	if user == nil {
		panic(exception.ApiException{
			Code:    []int{401},
			Message: "user cannot found",
		})
	}

	expiresAt := time.Now()
	if config.AuthConfig.Keep != 0 {
		expiresAt = expiresAt.Add(config.AuthConfig.Keep)
	} else {
		expiresAt = expiresAt.Add(time.Hour * 24 * 365 * 100)
	}

	userPK := prop.Get(config.AuthConfig.Model, config.AuthConfig.PrimaryKey)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    fmt.Sprintf("%v", userPK),
		ExpiresAt: expiresAt.Unix(),
	})

	token, err := claims.SignedString([]byte(config.AuthConfig.SecretKey))

	if err != nil {
		panic(exception.ApiException{
			Code:    []int{1002},
			Message: "cannot login",
		})
	}
	res := make(map[string]interface{})
	res["type"] = "Bearer"
	res["token"] = token
	res["expires_at"] = expiresAt.UTC()
	return res
}
