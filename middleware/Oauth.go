package middleware

import (
	"net/http"
	"p4_web/constant"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
)

func Oauth(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			// http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			// c.AbortWithError(http.StatusUnauthorized, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": []int{401},
				"data": err.Error(),
			})
		}
		c.Set(constant.OAUTH, info)
		c.Next()
	}
}
