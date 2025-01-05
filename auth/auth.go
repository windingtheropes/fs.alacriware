package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	// "fmt"

	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/windingtheropes/fs.alacriware/based"
)
type Credentials struct {
	UID int
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

// Check if path is within scope
func IsInPathScope(path string, scope string) bool {
	if strings.Count(path, "/") > strings.Count(scope, "/") {
		// The path is deeper than scope but not neccesarily within
		scope_dir := scope + "/"
		relative_path := strings.Replace(path, scope_dir, "", 1)
		if len(relative_path) + len(scope_dir) == len(path) {
			return true
		}
	}
	return false
}

// Whitelist basis; block by default
func canAccessResource(resource string, groups []int) bool {
	allowed := false
	if len(groups) == 0 {
		// User doesn't have any permissions
		return false
	} else {
		for i := range groups {
			gid := groups[i]
			permissions, err := based.DB.GetPermissions(gid)
			if err != nil {
				fmt.Println("Error getting group permissions: %v", err)
			}

			for p := range permissions {
				perm := permissions[p]
				fmt.Println(perm)
				if (perm.Apply_Recursive == true && IsInPathScope(resource, perm.Resource_Path)) || perm.Resource_Path == resource {
					if perm.Allowed == true { 
						// Soft alow
						allowed = true
						continue 
					} else { 
						// Hard deny
						return false 
					}
				} 
			}
		}
	}
	return allowed
}
// Authenticates on a whitelist basis, where all unauthenticated users are tried as User 1
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {		
		// if permitted
		credentials := getCredentials(c.Query("t"))
		groups, err := based.DB.GetUserMembership(credentials.UID)
		fmt.Println(c.Request.URL.Path)
		if err != nil {
			fmt.Println("Error getting user membership: %v", err)
		}
		if canAccessResource(c.Request.URL.Path, groups) == false {
			c.AbortWithStatus(403)
			return
		} else {
			c.Next()
		}
	}
}


