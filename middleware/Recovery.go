package middleware

import (
	"net/http"
	"p4_web/tools/exception"

	"github.com/gin-gonic/gin"
)

func MyRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// log.Print(recovered)
		if err, ok := recovered.(exception.ApiException); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": err.Code,
				"data": err.Message,
			})
			// panic("")
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
