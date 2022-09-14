package flow

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/model"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GeneralAllSources(c *gin.Context) {
	c.JSON(http.StatusOK, commom.SuccessPayload(model.GloOther.IDC))
}

func FlowTplPostSourceTemplate(c *gin.Context) {
	req := new(flowReq)
	if err := c.Bind(req); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	var ft model.CoreWorkflowTpl
	step, _ := json.Marshal(req.Steps)
	if model.DB().Where("source =?", req.Source).First(&ft).RecordNotFound() {
		model.DB().Create(&model.CoreWorkflowTpl{Source: req.Source, Steps: step})
	} else {
		model.DB().Model(model.CoreWorkflowTpl{}).Where("source =?", req.Source).Update(model.CoreWorkflowTpl{Steps: step})
	}

	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.DATA_IS_UPDATED))
}

func FlowTplEditSourceTemplateInfo(c *gin.Context) {
	req := new(flowReq)
	if err := c.Bind(req); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	if source := c.Query("source"); "" != source {
		req.Source = source
	}
	var ft model.CoreWorkflowTpl
	if model.DB().Where("source =?", req.Source).First(&ft).RecordNotFound() {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(errors.New("环境没有添加流程!无法审批工单")))
		return
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(ft))
}
