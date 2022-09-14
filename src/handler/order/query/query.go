package query

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func FetchQueryRecord(c *gin.Context) {
	u := new(commom.PageInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	order := u.GetSQLQueryList(
		commom.AccordingToQueryPer(),
		commom.AccordingToWorkId(u.Find.Text),
		commom.AccordingToDate(u.Find.Picker),
	)
	c.JSON(http.StatusOK, commom.SuccessPayload(order))
}

func FetchQueryOrder(c *gin.Context) {

	u := new(commom.PageInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)
	order := u.GetSQLQueryList(
		commom.AccordingToUsername(u.Find.Text),
		commom.AccordingToAssigned(user),
		commom.AccordingToDate(u.Find.Picker),
		commom.AccordingToAllQueryOrderState(u.Find.Status),
	)
	c.JSON(http.StatusOK, commom.SuccessPayload(order))
}

func FetchQueryRecordProfile(c *gin.Context) {
	u := new(commom.ExecuteStr)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	start, end := lib.Paging(u.Page, 20)
	var detail []model.CoreQueryRecord
	var count int
	model.DB().Model(&model.CoreQueryRecord{}).Where("work_id =?", u.WorkId).Count(&count).Offset(start).Limit(end).Find(&detail)
	c.JSON(http.StatusOK, commom.SuccessPayload(commom.CommonList{Data: detail, Page: count}))
}

func QueryDeleteEmptyRecord(c *gin.Context) {
	var j []model.CoreQueryOrder
	model.DB().Select("work_id").Where(`query_per =?`, 3).Find(&j)
	for _, i := range j {
		var k model.CoreQueryRecord
		if model.DB().Where("work_id =?", i.WorkId).First(&k).RecordNotFound() {
			model.DB().Where("work_id =?", i.WorkId).Delete(&model.CoreQueryOrder{})
		}
	}
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_CLEAR))
}

func QueryHandlerSets(c *gin.Context) {
	u := new(commom.QueryOrder)
	var s model.CoreQueryOrder
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	found := !model.DB().Where("work_id=? AND query_per=?", u.WorkId, 2).First(&s).RecordNotFound()
	switch u.Tp {
	case "agreed":
		if found {
			model.DB().Model(model.CoreQueryOrder{}).Where("work_id =?", u.WorkId).Update(map[string]interface{}{"query_per": 1, "ex_date": time.Now().Format("2006-01-02 15:04")})
			lib.MessagePush(u.WorkId, lib.EVENT_ORDER_QUERY_PASS, "")
		}
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_AGREE))
	case "reject":
		if found {
			model.DB().Model(model.CoreQueryOrder{}).Where("work_id =?", u.WorkId).Update(map[string]interface{}{"query_per": 0})
			lib.MessagePush(u.WorkId, lib.EVENT_ORDER_QUERY_REJECT, "")
		}
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_REJECT))
	case "stop":
		model.DB().Model(model.CoreQueryOrder{}).Where("work_id =?", u.WorkId).Update(map[string]interface{}{"query_per": 3})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_ALL_END))
	case "cancel":
		model.DB().Model(model.CoreQueryOrder{}).Updates(&model.CoreQueryOrder{QueryPer: 3})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_ALL_CANCEL))
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func AuditOrRecordQueryOrderFetchApis(c *gin.Context) {
	switch c.Param("tp") {
	case "list":
		FetchQueryOrder(c)
	case "record":
		FetchQueryRecord(c)
	case "profile":
		FetchQueryRecordProfile(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}
