package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// JwtConfig defines the config of JWT middleware
type JwtConfig struct {
	GetKey         string
	AuthScheme     string
	SigningKey     interface{}
	SigningMethod  string
	TokenLookup    string
	Claims         jwt.Claims
	keyFunc        jwt.Keyfunc
	ErrorHandler   JWTErrorHandler
	SuccessHandler JWTSuccessHandler
}

type jwtExtractor func(*gin.Context) (string, error)

// JWTErrorHandler defines a function which is error for a valid token.
type JWTErrorHandler func(error) error

// JWTSuccessHandler defines a function which is executed for a valid token.
type JWTSuccessHandler func(*gin.Context)

const algorithmHS256 = "HS256"

// DefaultJwtConfig is the default config of JWT middleware
var DefaultJwtConfig = JwtConfig{
	GetKey:        "auth",
	SigningMethod: algorithmHS256,
	AuthScheme:    "Bearer",
	TokenLookup:   "header:" + HeaderAuthorization,
	Claims:        jwt.MapClaims{},
}

// JWTWithConfig is the custom implementation CORS middleware
func JWTWithConfig(config JwtConfig) gin.HandlerFunc {
	if config.SigningKey == nil {
		panic("jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJwtConfig.SigningMethod
	}
	if config.GetKey == "" {
		config.GetKey = DefaultJwtConfig.GetKey
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJwtConfig.AuthScheme
	}

	if config.Claims == nil {
		config.Claims = DefaultJwtConfig.Claims
	}

	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJwtConfig.TokenLookup
	}

	config.keyFunc = func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return config.SigningKey, nil
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], config.AuthScheme)

	return func(c *gin.Context) {
		auth, err := extractor(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		token := new(jwt.Token)
		if _, ok := config.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, config.keyFunc)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				return
			}
		} else {
			t := reflect.ValueOf(config.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
		}
		if err == nil && token.Valid {
			c.Set(config.GetKey, token)
			return
		}
		// bug fix
		// if  invalid or expired jwt,
		// we must intercept all handlers and return serverError
		c.JSON(http.StatusUnauthorized, "invalid or expired jwt")
		return
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		auth := c.GetHeader(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or malformed jwt")
	}
}
