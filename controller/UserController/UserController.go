package UserController

import (
	"p4_web/constant"
	db "p4_web/database"
	"p4_web/model"
	"p4_web/tools/env"
	"p4_web/tools/exception"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetList(c *gin.Context) {
	var users []model.User
	db.DB.Find(&users)
	c.Set(constant.RESPONSE, users)
}

func Register(c *gin.Context) {
	var req map[string]string
	c.BindJSON(&req)

	password, _ := bcrypt.GenerateFromPassword([]byte(req["password"]), 14)

	user := model.User{
		Name:     req["name"],
		Mail:     req["mail"],
		Password: password,
	}
	db.DB.Create(&user)
	c.Set(constant.RESPONSE, user)
}

func Login(c *gin.Context) {
	var req map[string]string
	c.BindJSON(&req)

	var data model.User
	db.DB.Where("mail = ?", req["mail"]).First(&data)

	if data.Id == 0 {
		panic(exception.ApiException{
			Code:    []int{1001},
			Message: "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(data.Password, []byte(req["password"])); err != nil {
		panic(exception.ApiException{
			Code:    []int{1002},
			Message: "password not equal",
		})
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(data.Id)),
		ExpiresAt: expiresAt.Unix(),
	})

	token, err := claims.SignedString([]byte(env.Get("SECRET_KEY")))

	if err != nil {
		panic(exception.ApiException{
			Code:    []int{1002},
			Message: "could not login",
		})
	}

	c.Set(constant.RESPONSE, gin.H{
		"type":       "Bearer",
		"token":      token,
		"expires_at": expiresAt.UTC(),
	})
}

func GetSelf(c *gin.Context) {
	// user, _ := c.Get(constant.AUTH)
	// c.Set(constant.RESPONSE, user)
}
