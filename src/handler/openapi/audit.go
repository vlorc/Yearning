package openapi

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/handler/manager/flow"
	"Yearning-go/src/handler/order/audit"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExecuteStr struct {
	WorkId   string `json:"work_id"`
	Perform  string `json:"perform"`
	Page     int    `json:"page"`
	Flag     int    `json:"flag"`
	Text     string `json:"text"`
	Tp       string `json:"tp"`
	Assigned string `json:"assigned"`
}

func AuditOrder(c *gin.Context) {
	u := new(ExecuteStr)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	// add user check
	var order model.CoreSqlOrder
	model.DB().Where("work_id =?", u.WorkId).First(&order)

	var flowTpl model.CoreWorkflowTpl
	model.DB().Where("source =?", order.IDC).First(&flowTpl)

	var steps []flow.Step
	_ = json.Unmarshal(flowTpl.Steps, &steps)

	if order.Status != 2 && order.Status != 5 || order.CurrentStep >= len(steps) {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("工单状态不允许审批")))
		return
	}

	if u.Assigned == order.Username {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("你是提交人不能进行操作")))
		return
	}

	var authCheck bool

	for _, v := range steps[order.CurrentStep].Auditor {
		if v == u.Assigned {
			authCheck = true
			break
		}
	}

	if !authCheck {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(fmt.Errorf("你没有权限进行操作")))
		return
	}

	req := &commom.ExecuteStr{
		WorkId:  u.WorkId,
		Perform: u.Perform,
		Page:    u.Page,
		Flag:    u.Flag,
		Text:    u.Text,
		Tp:      u.Tp,
	}
	if "" == u.Perform {
		u.Perform = u.Assigned
	}
	if 0 == req.Flag {
		req.Flag = order.CurrentStep
	}

	switch u.Tp {
	case "agree":
		if order.CurrentStep+1 == len(steps) {
			c.JSON(http.StatusOK, audit.ExecuteOrderEx(req, u.Assigned))
		} else {
			auditor := steps[order.CurrentStep+1].Auditor
			if len(auditor) > 0 && stringsIndexOf(auditor, u.Perform) < 0 {
				u.Perform = steps[order.CurrentStep].Auditor[0]
			}
			c.JSON(http.StatusOK, audit.MultiAuditOrder(req, u.Assigned))
		}
	case "reject":
		c.JSON(http.StatusOK, audit.RejectOrder(req, u.Assigned))
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func stringsIndexOf(strs []string, val string) int {
	for i, s := range strs {
		if val == s {
			return i
		}
	}
	return -1
}
