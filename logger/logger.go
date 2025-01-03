package logger

import (	
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/windingtheropes/fs.alacriware/based"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {		
		c.Next()
		// after request
		// don't log nonexistant routes
		
		if(c.Writer.Status() != 404) {
			req := based.Request_Log {
				IP: c.ClientIP(),
				Access_Time: time.Now().UnixMilli(),
				Resource_Path: c.Request.RequestURI,
				Token: c.Query("t"),
				Code: c.Writer.Status(),
			}
			_, err := based.DB.LogRequest(req)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
