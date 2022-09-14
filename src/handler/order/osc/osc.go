package osc

import (
	"Yearning-go/src/handler/commom"
	"github.com/gin-gonic/gin"
	"net/http"
)

// OscPercent show OSC percent
func OscPercent(c *gin.Context) {
	var k = &OSC{WorkId: c.Param("work_id")}
	c.JSON(http.StatusOK, commom.SuccessPayload(k.Percent()))
}

// OscKill will kill OSC command
func OscKill(c *gin.Context) {
	var k = OSC{WorkId: c.Param("work_id")}
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(k.Kill()))
}
