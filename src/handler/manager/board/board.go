package board

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/model"
	"net/http"
	"github.com/gin-gonic/gin"
)

type boardReq struct {
	Board string `json:"board"`
}

const BOARD_MESSAGE_SAVE = "公告已保存"

func GeneralPostBoard(c *gin.Context) {
	req := new(boardReq)
	if err := c.Bind(req); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	model.DB().Model(model.CoreGlobalConfiguration{}).Update(&model.CoreGlobalConfiguration{Board: req.Board})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(BOARD_MESSAGE_SAVE))
}
