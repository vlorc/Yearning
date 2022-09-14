package user

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuperFetchUser(c *gin.Context) {
	u := new(fetchUser)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	var user []model.CoreAccount
	var count int
	start, end := lib.Paging(u.Page, 10)
	if u.Find.Valve {
		model.DB().Model(model.CoreAccount{}).Select(CommonExpr).Scopes(
			commom.AccordingToUsername(u.Find.Username),
			commom.AccordingToOrderDept(u.Find.Dept),
		).Count(&count).Offset(start).Limit(end).Find(&user)
	} else {
		model.DB().Model(model.CoreAccount{}).Select(CommonExpr).Count(&count).Offset(start).Limit(end).Find(&user)
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(commom.CommonList{
		Page:  count,
		Data:  user,
		Multi: model.GloOther.Multi,
	}))
}

func SuperDeleteUser(c *gin.Context) {
	user := c.Query("user")
	if user == "admin" {
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(ADMIN_NOT_DELETE))
		return
	}
	tx := model.DB().Begin()
	model.DB().Where("username =?", user).Delete(&model.CoreAccount{})
	model.DB().Where("username =?", user).Delete(&model.CoreGrained{})
	tx.Commit()

	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(USER_DELETE_SUCCESS, user)))
}

func ManageUserCreateOrEdit(c *gin.Context) {
	u := new(CommonUserPost)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	switch u.Tp {
	case "edit":
		c.JSON(http.StatusOK, SuperUserEdit(&u.User))
	case "create":
		c.JSON(http.StatusOK, SuperUserRegister(&u.User))
	case "password":
		model.DB().Model(&model.CoreAccount{}).Where("username = ?", u.User.Username).Update(
			model.CoreAccount{Password: lib.DjangoEncrypt(u.User.Password, string(lib.GetRandom()))})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(USER_EDIT_PASSWORD_SUCCESS))
	default:
		c.JSON(http.StatusOK,commom.ERR_REQ_FAKE)
	}
}

func ManageUserFetch(c *gin.Context) {
	u := new(CommonUserGet)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	switch u.Tp {
	case "depend":
		c.JSON(http.StatusOK, DelUserDepend(u.User))
	case "group":
		var p []model.CoreRoleGroup
		var userP model.CoreGrained
		model.DB().Find(&p)
		model.DB().Where("username=?", u.User).First(&userP)
		c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"group": userP.Group, "list": p}))
	default:
		c.JSON(http.StatusOK,commom.ERR_REQ_FAKE)
	}
}
