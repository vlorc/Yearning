package personal

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PersonalFetchMyOrder(c *gin.Context){
	u := new(commom.PageInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)

	var pg int

	var order []model.CoreSqlOrder

	start, end := lib.Paging(u.Page, 15)

	model.DB().Model(&model.CoreSqlOrder{}).Select(commom.QueryField).
		Scopes(
			commom.AccordingToAllOrderState(u.Find.Status),
			commom.AccordingToUsernameEqual(user),
			commom.AccordingToDatetime(u.Find.Picker),
			commom.AccordingToText(u.Find.Text),
		).Order("id desc").Count(&pg).Offset(start).Limit(end).Find(&order)
	c.JSON(http.StatusOK, commom.SuccessPayload(commom.CommonList{Data: order, Page: pg, Multi: model.GloOther.Multi}))
	return
}

func PersonalUserEdit(c *gin.Context) {
	param := c.Query("tp")
	u := new(model.CoreAccount)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	user, _ := lib.JwtParse(c)
	switch param {
	case "password":
		model.DB().Model(&model.CoreAccount{}).Where("username = ?", user).Update(
			&model.CoreAccount{Password: lib.DjangoEncrypt(u.Password, string(lib.GetRandom()))})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(CUSTOM_PASSWORD_SUCCESS))
	case "mail":
		model.DB().Model(&model.CoreAccount{}).Where("username = ?", user).Updates(model.CoreAccount{Email: u.Email, RealName: u.RealName, Mobile: u.Mobile})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(CUSTOM_INFO_SUCCESS))
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}

func PersonalFetchOrderListOrProfile(c *gin.Context) {
	switch c.Param("tp") {
	case "list":
		PersonalFetchMyOrder(c)
	case "edit":
		PersonalUserEdit(c)
	default:
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
}
