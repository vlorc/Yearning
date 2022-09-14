package audit

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	pb "Yearning-go/src/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SuperSQLTest(c *gin.Context) {
	u := new(commom.SQLTest)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	var s model.CoreDataSource
	var order model.CoreSqlOrder
	model.DB().Where("work_id =?", u.WorkId).First(&order)
	model.DB().Where("source =?", order.Source).First(&s)
	y := pb.LibraAuditOrder{
		IsDML:    order.Type == 1,
		SQL:      order.SQL,
		DataBase: order.DataBase,
		Source: &pb.Source{
			Addr:     s.IP,
			User:     s.Username,
			Port:     int32(s.Port),
			Password: lib.Decrypt(s.Password),
		},
		Execute: false,
		Check:   true,
	}
	record, err := lib.TsClient(&y, s.Proxy)
	if err != nil {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		return
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(record))
}

func ExecuteOrder(c *gin.Context) {
	u := new(commom.ExecuteStr)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)
	var order model.CoreSqlOrder

	model.DB().Where("work_id =?", u.WorkId).First(&order)

	if order.Status != 2 && order.Status != 5 {
		// c.Logger().Error(IDEMPOTENT)
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(IDEMPOTENT))
		return
	}

	if order.Type == 3 {
		model.DB().Model(&model.CoreSqlOrder{}).Where("work_id =?", u.WorkId).Updates(map[string]interface{}{"status": 1, "execute_time": time.Now().Format("2006-01-02 15:04"), "current_step": order.CurrentStep + 1})
	} else {
		executor := new(Review)

		order.Assigned = user

		executor.Init(order).Executor()
	}
	model.DB().Create(&model.CoreWorkflowDetail{
		WorkId:   u.WorkId,
		Username: user,
		Rejected: "",
		Time:     time.Now().Format("2006-01-02 15:04"),
		Action:   ORDER_EXECUTE_STATE,
	})
	
	lib.MessagePush(u.WorkId, lib.EVENT_ORDER_EXEC_PASS, "")

	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(ORDER_EXECUTE_STATE))
}

func AuditOrderState(c *gin.Context) {
	u := new(commom.ExecuteStr)
	user, _ := lib.JwtParse(c)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	switch u.Tp {
	case "agree":
		c.JSON(http.StatusOK, MultiAuditOrder(u, user))
	case "reject":
		c.JSON(http.StatusOK, RejectOrder(u, user))
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

// DelayKill will stop delay order
func DelayKill(c *gin.Context) {
	u := new(commom.ExecuteStr)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)
	model.DB().Create(&model.CoreWorkflowDetail{
		WorkId:   u.WorkId,
		Username: user,
		Time:     time.Now().Format("2006-01-02 15:04"),
		Action:   ORDER_KILL_STATE,
	})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(delayKill(u.WorkId)))
}

func FetchAuditOrder(c *gin.Context) {
	u := new(commom.PageInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)
	order := u.GetSQLOrderList(commom.AccordingToAllOrderState(u.Find.Status),
		commom.AccordingToRelevant(user),
		commom.AccordingToText(u.Find.Text),
		commom.AccordingToDatetime(u.Find.Picker))
	c.JSON(http.StatusOK, commom.SuccessPayload(order))
}

func FetchRecord(c *gin.Context) {
	u := new(commom.PageInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	order := u.GetSQLOrderList(commom.AccordingToOrderState(),
		commom.AccordingToWorkId(u.Find.Text),
		commom.AccordingToDatetime(u.Find.Picker))
	c.JSON(http.StatusOK, commom.SuccessPayload(order))
}

func AuditOrderApis(c *gin.Context) {
	switch c.Param("tp") {
	case "test":
		SuperSQLTest(c)
	case "state":
		AuditOrderState(c)
	case "execute":
		ExecuteOrder(c)
	case "kill":
		DelayKill(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func AuditOrRecordOrderFetchApis(c *gin.Context) {
	switch c.Param("tp") {
	case "list":
		FetchAuditOrder(c)
	case "record":
		FetchRecord(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}
