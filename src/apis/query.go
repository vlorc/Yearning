package apis

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/handler/fetch"
	"Yearning-go/src/handler/personal"
	"Yearning-go/src/lib"
	"github.com/gin-gonic/gin"
	"net/http"
)

func YearningQueryForGet(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "fetch_table":
		personal.FetchQueryTableInfo(c)
	case "table_info":
		personal.FetchQueryTableStruct(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningQueryForPut(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "fetch_base":
		personal.FetchQueryDatabaseInfo(c)
	case "status":
		personal.FetchQueryStatus(c)
	case "merge":
		fetch.FetchMergeDDL(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningQueryForPost(c *gin.Context) {
	tp := c.Param("tp")
	user, _ := lib.JwtParse(c)
	switch tp {
	case "refer":
		personal.ReferQueryOrder(c, &user)
	case "results":
		personal.FetchQueryResults(c, &user)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningQueryApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:    YearningQueryForGet,
		Put:    YearningQueryForPut,
		Post:   YearningQueryForPost,
		Delete: personal.UndoQueryOrder,
	}
}
