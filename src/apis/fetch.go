package apis

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/handler/fetch"
	"Yearning-go/src/handler/manager/group"
	"Yearning-go/src/lib"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FetchResourceForGet(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "detail":
		fetch.FetchOrderDetailList(c)
	case "roll":
		fetch.FetchOrderDetailRollSQL(c)
	case "undo":
		fetch.FetchUndo(c)
	case "sql":
		fetch.FetchSQLInfo(c)
	case "perform":
		fetch.FetchPerformList(c)
	case "idc":
		fetch.FetchIDC(c)
	case "source":
		fetch.FetchSource(c)
	case "base":
		fetch.FetchBase(c)
	case "table":
		fetch.FetchTable(c)
	case "fields":
		fetch.FetchTableInfo(c)
	case "steps":
		fetch.FetchStepsProfile(c)
	case "board":
		fetch.FetchBoard(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func FetchResourceForPut(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "test":
		fetch.FetchSQLTest(c)
	case "merge":
		fetch.FetchMergeDDL(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func FetchResourceForPost(c *gin.Context) {
	tp := c.Param("tp")
	switch tp {
	case "marge":
		group.SuperUserRuleMarge(c)
	case "roll_order":
		fetch.RollBackSQLOrder(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func YearningFetchApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Post: FetchResourceForPost,
		Get:  FetchResourceForGet,
		Put:  FetchResourceForPut,
	}
}
