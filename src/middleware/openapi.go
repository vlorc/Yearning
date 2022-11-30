package middleware

import (
	"Yearning-go/src/handler/commom"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/util/math"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewOpenApi(appid, secret string) func(*gin.Context) {
	if "" == appid {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("openapi not enable")))
			c.Abort()
			return
		}
	}
	if "" == secret {
		return NewTokenOpenApi(appid)
	}

	appidBytes := []byte(appid)
	secretBytes := []byte(secret)
	return func(c *gin.Context) {
		token := c.Query("token")
		values := strings.Split(token, ".")
		if len(values) < 3 {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("token not exist")))
			c.Abort()
			return
		}
		if values[0] != appid {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("token appid not exist")))
			c.Abort()
			return
		}
		if ts, _ := strconv.ParseInt(values[1], 10, 63); math.Abs(time.Now().Unix()-ts) > 60 {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("token timestamp greater than 60")))
			c.Abort()
			return
		}

		data, _ := io.ReadAll(c.Request.Body)
		if len(data) == 0 {
			return
		}

		h := hmac.New(md5.New, secretBytes)

		h.Write(appidBytes)
		h.Write([]byte(values[1]))
		h.Write(data)

		b := h.Sum(nil)
		sign := hex.EncodeToString(b)

		if sign != values[2] {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("token verification failed")))
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(data))
	}
}

func NewTokenOpenApi(appid string) func(*gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")
		values := strings.Split(token, ".")
		if len(values) < 1 || values[0] != appid {
			c.Error(errors.New("token not exist"))
			c.Abort()
			return
		}
	}
}
