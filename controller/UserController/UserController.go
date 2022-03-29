package UserController

import (
	"fmt"
	"p4_web/constant"
	db "p4_web/database"
	"p4_web/model"
	"p4_web/tools/auth"
	"p4_web/tools/exception"

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

	user := model.User{
		Name:     req["name"],
		Mail:     req["mail"],
		Password: []byte(req["password"]),
	}
	db.DB.Create(&user)
	c.Set(constant.RESPONSE, user)
}

func Login(c *gin.Context) {
	var req map[string]string
	c.BindJSON(&req)

	var user model.User
	err := db.DB.Where("mail = ?", req["mail"]).First(&user).Error
	if err != nil {
		panic(exception.ApiException{
			Code:    []int{1001},
			Message: fmt.Sprintf("user not found: %v", err.Error()),
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req["password"])); err != nil {
		panic(exception.ApiException{
			Code:    []int{1002},
			Message: fmt.Sprintf("password not equal: %v", err.Error()),
		})
	}

	c.Set(constant.RESPONSE, auth.Generate())
}

func GetSelf(c *gin.Context) {
	user, _ := c.Get(constant.AUTH)
	c.Set(constant.RESPONSE, user)
}
