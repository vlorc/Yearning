package personal

import (
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	model.DbInit("../../../conf.toml")
}

func QueryRes(c *gin.Context) {
	user := "admin"
	FetchQueryResults(c, &user)
}

func BenchmarkFetchQueryResults(b *testing.B) {
	c := gin.New()
	c.POST("/api/v2/query", QueryRes)
	b.ReportAllocs()
	b.SetBytes(1024 * 1024)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v2/query", strings.NewReader(`{"sql":"select * from core_accounts","data_base":"public","source":"local"}`))
		req.Header.Set("Content-Type", gin.MIMEJSON)
		rec := httptest.NewRecorder()
		c.ServeHTTP(rec, req)
	}
}
