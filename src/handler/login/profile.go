package login

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserReqSwitch(c *gin.Context) {
	c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"reg": model.GloOther.Register, "valid": true}))
}
