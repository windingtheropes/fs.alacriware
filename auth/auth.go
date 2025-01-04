package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	// "fmt"

	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/windingtheropes/fs.alacriware/based"
)
type Credentials struct {
	UID int
	Token string
}
const TOKEN_LENGTH  = 64

func NewToken(uid int, expiry int, max_uses int) string {
	randomBytes := make([]byte, TOKEN_LENGTH)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
	return base64.URLEncoding.EncodeToString(randomBytes)[:TOKEN_LENGTH]
}
func getCredentials(tQuery string) Credentials {
		res, err := based.DB.FindToken(tQuery)
		if err != nil {
			fmt.Println(err)
		}
		if len(res) < 1 {
			// Token doesn't exist
			return Credentials{
				UID: 1,
			}
		} else {
			// Token exists, make sure it is valid before authing as the user contained
			token := res[0]
			if(token.Used > token.Max) {
				// Token is invalid due to max uses
				return Credentials{
					UID: 1,
				}
			}
			if(time.Now().UnixMilli() > token.Expiry) {
				// Token is invalid due to ttl
				return Credentials{
					UID: 1,
				}
			}
			// The token is valid, so the user is acceptable
			return Credentials{
				UID: token.User_ID,
			}
		}
}
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {		
		// if permitted
		c.Next()
		credentials := getCredentials(c.Query("t"))
		fmt.Println(credentials)
		// else
		// c.Status(403)
		// c.Abort()

	}
}

