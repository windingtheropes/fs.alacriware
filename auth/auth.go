package auth

import (
	"crypto/rand"
	"encoding/base64"
	// "fmt"

	"github.com/gin-gonic/gin"
	// "github.com/windingtheropes/fs.alacriware/based"
)

const TOKEN_LENGTH  = 64

func NewToken(uid int, expiry int, max_uses int) string {
	randomBytes := make([]byte, TOKEN_LENGTH)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
	return base64.URLEncoding.EncodeToString(randomBytes)[:TOKEN_LENGTH]
}
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {		
		// if permitted
		c.Next()
		// else
		// c.Status(403)
		// c.Abort()

	}
}

