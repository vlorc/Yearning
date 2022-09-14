package roles

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	ser "Yearning-go/src/parser"
	pb "Yearning-go/src/proto"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuperSaveRoles(c *gin.Context) {

	u := new(ser.AuditRole)

	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	ser.FetchAuditRole = *u
	audit, _ := json.Marshal(u)
	model.DB().Model(model.CoreGlobalConfiguration{}).Updates(&model.CoreGlobalConfiguration{AuditRole: audit})
	lib.OverrideConfig(&pb.LibraAuditOrder{})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.DATA_IS_EDIT))
}

func SuperFetchRoles(c *gin.Context) {
	var k model.CoreGlobalConfiguration
	model.DB().Select("audit_role").First(&k)
	c.JSON(http.StatusOK, commom.SuccessPayload(k))
}
