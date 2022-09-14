package apis

import (
	"Yearning-go/src/handler"
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"github.com/gin-gonic/gin"
	"net/http"
)

func YearningDashGet(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "pie":
		handler.DashPie(c)
	case "axis":
		handler.DashAxis(c)
	case "count":
		handler.DashCount(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningDashPut(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "profile":
		handler.DashUserInfo(c)
	case "stmt":
		handler.DashStmt(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningDashApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get: YearningDashGet,
		Put: YearningDashPut,
	}
}
