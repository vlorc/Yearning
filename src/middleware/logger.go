package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf("http request method: %s path: %s client: %s latency: %v status: %d\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			time.Now().Sub(start),
			c.Writer.Status(),
		)
	}
}

