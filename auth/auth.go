package auth

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if permitted
		c.Next()
		// else
		// c.Status(403)
		// c.Abort()

	}
}
