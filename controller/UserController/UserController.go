package UserController

import (
	"p4_web/constant"
	db "p4_web/database"
	"p4_web/model"
	"p4_web/tools/auth"

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
	auth.Check(req["mail"], req["password"])

	c.Set(constant.RESPONSE, auth.Generate())
}

func GetSelf(c *gin.Context) {
	user, _ := c.Get(constant.AUTH)
	c.Set(constant.RESPONSE, user)
}
