package auth

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"p4_web/config"

	"p4_web/constant"
	db "p4_web/database"
	"p4_web/tools/exception"
	"p4_web/tools/prop"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware(guards ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var guarder *config.Guard
		var err error
		if len(guards) == 0 {
			guards = append(guards, config.AuthConfig.Guard)
		}
		for _, guard := range guards {
			if guarder == nil {
				guarder, err = Check(c, guard)
			}
		}
		if guarder == nil {
			panic(exception.ApiException{
				Code:    []int{http.StatusUnauthorized},
				Message: err.Error(),
			})
		}
		c.Set(constant.AUTH, guarder.Model)
		c.Next()
	}
}

func Check(c *gin.Context, guard string) (*config.Guard, error) {
	s := strings.Split(c.Request.Header["Authorization"][0], " ")
	if s[0] == "Bearer" {
		guarder := config.AuthConfig.Guarder(guard)
		token, err := jwt.ParseWithClaims(s[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(guarder.SecretKey), nil
		})
		if err != nil {
			return nil, err
		}
		claims := token.Claims.(*jwt.StandardClaims)
		s = strings.Split(claims.Issuer, ";")
		userPK, tokenGuard := s[0], s[1]
		if tokenGuard != guard {
			return nil, errors.New("guard error")
		}
		err = db.DB.Where(fmt.Sprintf("%v = ?", guarder.PrimaryKey), userPK).First(&guarder.Model).Error
		if err != nil {
			return nil, err
		}
		return &guarder, nil
	} else {
		return nil, errors.New("invalid access token")
	}
}

func Generate(guard string, checkUser interface{}) map[string]interface{} {
	guarder := config.AuthConfig.Guarder(guard)
	// user := guarder.Model
	if checkUser == nil {
		panic(exception.ApiException{
			Code:    []int{401},
			Message: "user cannot found",
		})
	}
	if reflect.TypeOf(checkUser) != reflect.TypeOf(guarder.Model) {
		panic(exception.ApiException{
			Code:    []int{401},
			Message: "user not equal",
		})
	}
	user := checkUser

	expiresAt := time.Now()
	if guarder.Keep != 0 {
		expiresAt = expiresAt.Add(guarder.Keep)
	} else {
		expiresAt = expiresAt.Add(time.Hour * 24 * 365 * 100)
	}
	userPK := prop.Get(user, guarder.PrimaryKey)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    fmt.Sprintf("%v;%v", userPK, guard),
		ExpiresAt: expiresAt.Unix(),
	})

	token, err := claims.SignedString([]byte(guarder.SecretKey))

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
