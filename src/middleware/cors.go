package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig defined the config of CORS middleware
type CORSConfig struct {
	Origins          []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

// DefaultCORSConfig is the default config of  CORS middleware
var DefaultCORSConfig = CORSConfig{
	Origins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
		http.MethodConnect,
		http.MethodTrace,
	},
}

// Cors is the default implementation CORS middleware
func Cors() gin.HandlerFunc {
	return CorsWithConfig(DefaultCORSConfig)
}

// CorsWithConfig is the default implementation CORS middleware
func CorsWithConfig(config CORSConfig) gin.HandlerFunc {

	if len(config.Origins) == 0 {
		config.Origins = DefaultCORSConfig.Origins
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")

	allowHeaders := strings.Join(config.AllowHeaders, ",")

	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	maxAge := strconv.Itoa(config.MaxAge)

	return func(c *gin.Context)  {

		localOrigin := c.GetHeader(HeaderOrigin)

		allowOrigin := ""

		m := c.Request.Method

		for _, o := range config.Origins {
			if o == "*" && config.AllowCredentials {
				allowOrigin = localOrigin
				break
			}
			if o == "*" || o == localOrigin {
				allowOrigin = o
				break
			}
		}

		// when method was not OPTIONS,
		// we can return simple response header
		// because the OPTIONS method is used to
		// describe the communication options for the target resource
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS

		if m != http.MethodOptions {
			c.Writer.Header().Add(HeaderVary, HeaderOrigin)
			c.Header(HeaderAccessControlAllowOrigin, allowOrigin)
			if config.AllowCredentials {
				c.Header(HeaderAccessControlAllowCredentials, "true")
			}
			if exposeHeaders != "" {
				c.Header(HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			c.Next()
			return
		}

		c.Writer.Header().Add(HeaderVary, HeaderOrigin)
		c.Writer.Header().Add(HeaderVary, HeaderAccessControlRequestMethod)
		c.Writer.Header().Add(HeaderVary, HeaderAccessControlRequestHeaders)
		c.Writer.Header().Set(HeaderAccessControlAllowOrigin, allowOrigin)
		c.Writer.Header().Set(HeaderAccessControlAllowMethods, allowMethods)
		if config.AllowCredentials {
			c.Writer.Header().Set(HeaderAccessControlAllowCredentials, "true")
		}
		if allowHeaders != "" {
			c.Writer.Header().Set(HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			h := c.GetHeader(HeaderAccessControlRequestHeaders)
			if h != "" {
				c.Writer.Header().Set(HeaderAccessControlAllowHeaders, h)
			}
		}
		if config.MaxAge > 0 {
			c.Writer.Header().Set(HeaderAccessControlMaxAge, maxAge)
		}
		return
	}
}
