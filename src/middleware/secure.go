package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Header types
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

type (

	//SecureConfig defines config of secure middleware
	SecureConfig struct {
		XSSProtection string `yaml:"xss_protection"`

		ContentTypeNosniff string `yaml:"content_type_nosniff"`

		XFrameOptions string `yaml:"x_frame_options"`

		HSTSMaxAge int `yaml:"hsts_max_age"`

		HSTSExcludeSubdomains bool `yaml:"hsts_exclude_subdomains"`

		ContentSecurityPolicy string `yaml:"content_security_policy"`

		CSPReportOnly bool `yaml:"csp_report_only"`

		HSTSPreloadEnabled bool `yaml:"hsts_preload_enabled"`

		ReferrerPolicy string `yaml:"referrer_policy"`
	}
)

// DefaultSecureConfig is default config of secure middleware
var DefaultSecureConfig = SecureConfig{
	XSSProtection:      "1; mode=block",
	ContentTypeNosniff: "nosniff",
	XFrameOptions:      "SAMEORIGIN",
	HSTSPreloadEnabled: false,
}

// Secure is default implementation of secure middleware
func Secure() gin.HandlerFunc {
	return SecureWithConfig(DefaultSecureConfig)
}

// SecureWithConfig is custom implementation of secure middleware
func SecureWithConfig(config SecureConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		if config.XSSProtection != "" {
			c.Header(HeaderXXSSProtection, config.XSSProtection)
		}

		if config.ContentTypeNosniff != "" {
			c.Header(HeaderXContentTypeOptions, config.ContentTypeNosniff)
		}

		if config.XFrameOptions != "" {
			c.Header(HeaderXFrameOptions, config.XFrameOptions)
		}

		if (c.GetHeader(HeaderXForwardedProto) == "https") && config.HSTSMaxAge != 0 {
			subdomains := ""
			if !config.HSTSExcludeSubdomains {
				subdomains = "; includeSubdomains"
			}
			if config.HSTSPreloadEnabled {
				subdomains = fmt.Sprintf("%s; preload", subdomains)
			}
			c.Header(HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", config.HSTSMaxAge, subdomains))
		}
		// CSP
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy-Report-Only
		// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Content_Security_Policy
		if config.ContentSecurityPolicy != "" {
			if config.CSPReportOnly {
				c.Header(HeaderContentSecurityPolicyReportOnly, config.ContentSecurityPolicy)
			} else {
				c.Header(HeaderContentSecurityPolicy, config.ContentSecurityPolicy)
			}
		}

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
		if config.ReferrerPolicy != "" {
			c.Header(HeaderReferrerPolicy, config.ReferrerPolicy)
		}
		return
	}
}
