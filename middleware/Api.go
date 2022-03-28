package middleware

import (
	"net/http"
	"p4_web/constant"

	"github.com/gin-gonic/gin"
)

func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		res, _ := c.Get(constant.RESPONSE)
		c.JSON(http.StatusOK, gin.H{
			"code": []int{0},
			"data": res,
		})
	}
}
