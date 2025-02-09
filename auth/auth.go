package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/windingtheropes/fs.alacriware/based/webdb"
)
type Credentials struct {
	UID int
}
const TOKEN_LENGTH  = 64

func genToken() string {
	randomBytes := make([]byte, TOKEN_LENGTH)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
	return base64.URLEncoding.EncodeToString(randomBytes)[:TOKEN_LENGTH]
}

func NewToken(uid int, expiry int64, max_uses int16) (webdb.Token, error) {
	var tok webdb.Token = webdb.Token{
		ID: genToken(),
		User_ID: uid,
		Expiry: expiry,
		Max: int16(max_uses),
		Used: 0,
	}
	_, err := webdb.WDB.AddToken(tok)
	return tok, err
}

func getCredentials(tQuery string) Credentials {
		res, err := webdb.WDB.FindToken(tQuery)
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
			// Check if the token is still alive, 0 means it won't expire
			if(time.Now().UnixMilli() > token.Expiry) && token.Expiry != 0 {
				// Token is invalid due to ttl
				return Credentials{
					UID: 1,
				}
			}

			// Increment uses on the token, and save to database
			// Only increment if within the limits, don't increment over
			// Take stats on unlimited tokens (when max is 0)
			if (token.Used <= token.Max) || token.Max == 0 {
				token.Used += 1
				_, err := webdb.WDB.UpdateToken(token)
				if err != nil {
					fmt.Printf("Token update error: %v", err)
					return Credentials {
						UID: 1,
					}
				}
			}
			
			// Token.Max == 0 means unlimited use
			if(token.Used > token.Max) && token.Max != 0 {
				// Token is invalid due to max uses
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
	// /hello , ensure never terminating slash
	if path != "/" { 
		path, _ = strings.CutSuffix(path, "/") 
	}
	if scope != "/" { 
		scope, _ = strings.CutSuffix(scope, "/") 
	}
	
	if strings.Count(path, "/") > strings.Count(scope, "/") {
		// The path is deeper than scope but not neccesarily within
		
		// Confirm most cases using the replace and add method, confirm by ensuring that the replaced portion was at the front of the path
		scope_explicit := func () string {
			if scope == "/" {
				return scope
			} else {
				return scope + "/"
			}
		}()
		
		relative_path := strings.Replace(path, scope_explicit, "", 1)
		if len(relative_path) + len(scope_explicit) == len(path) && strings.HasPrefix(path, scope) == true {
			return true
		}
	} else if scope == "/" {
		return true
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
			permissions, err := webdb.WDB.GetPermissions(gid)
			if err != nil {
				fmt.Printf("Error getting group permissions: %v", err)
			}

			for p := range permissions {
				perm := permissions[p]
				if (perm.Apply_Recursive && IsInPathScope(resource, perm.Resource_Path)) || perm.Resource_Path == resource {
					if perm.Allowed { 
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
		groups, err := webdb.WDB.GetUserMembership(credentials.UID)
		if err != nil {
			fmt.Printf("Error getting user membership: %v", err)
		}
		if !canAccessResource(c.Request.URL.Path, groups) {
			c.AbortWithStatus(403)
			return
		} else {
			c.Next()
		}
	}
}


