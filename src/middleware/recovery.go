package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func Recovery(dump func(int) []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				log.Printf("http crash method: %s path: %s client: %s error: %v\n", c.Request.Method, c.Request.URL.Path, c.ClientIP(), err)

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		c.Next()
	}
}
