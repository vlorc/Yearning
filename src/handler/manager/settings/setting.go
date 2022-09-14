// Copyright 2019 HenryYee.
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package settings

import (
	"Yearning-go/internal/pkg/messagex"
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	pb "Yearning-go/src/proto"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	WEBHOOK_TEST      = "测试消息已发送！请注意查收！"
	MAIL_TEST         = "测试邮件已发送！请注意查收！"
	ERR_LDAP_TEST     = "ldap连接失败!"
	SUCCESS_LDAP_TEST = "ldap连接成功!"
)

type settingInfo struct {
	Ldap    model.Ldap    `json:"ldap"`
	Message model.Message `json:"message"`
	Other   model.Other   `json:"other"`
}

type ber struct {
	Date string `json:"date"`
	Tp   bool   `json:"tp"`
}

func SuperFetchSetting(c *gin.Context) {

	var k model.CoreGlobalConfiguration

	model.DB().Select("ldap,message,other").First(&k)

	c.JSON(http.StatusOK, commom.SuccessPayload(k))
}

func SuperSaveSetting(c *gin.Context) {

	u := new(settingInfo)

	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	other, _ := json.Marshal(u.Other)
	message, _ := json.Marshal(u.Message)
	ldap, _ := json.Marshal(u.Ldap)
	diffIDC(u.Other.IDC)
	model.DB().Model(model.CoreGlobalConfiguration{}).Updates(&model.CoreGlobalConfiguration{Other: other, Message: message, Ldap: ldap})
	model.GloOther = u.Other
	model.GloLdap = u.Ldap
	model.GloMessage = u.Message
	lib.OverrideConfig(&pb.LibraAuditOrder{})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.DATA_IS_EDIT))
}

func diffIDC(src []string) {
	var idc model.CoreGlobalConfiguration
	var env model.Other
	model.DB().Find(&idc)
	_ = json.Unmarshal(idc.Other, &env)
	p := lib.NonIntersect(src, env.IDC)
	for _, i := range p {
		model.DB().Model(model.CoreWorkflowTpl{}).Where("source =?", i).Delete(&model.CoreWorkflowTpl{})
	}
}

func SuperTestSetting(c *gin.Context) {

	el := c.Query("test")
	u := new(settingInfo)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	if el == "mail" {
		go lib.SendMail(u.Message, messagex.Message{
			Body:   "test",
			Target: messagex.Target{Mobiles: []string{u.Message.ToUser}},
		})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(MAIL_TEST))
		return
	}

	if el == "web_hook" {
		go lib.WebHookTestPush(u.Message)
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(WEBHOOK_TEST))
		return
	}

	if el == "ldap" {
		if err := lib.LdapTest(&u.Ldap); nil == err {
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(SUCCESS_LDAP_TEST))
		} else {
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(ERR_LDAP_TEST))
		}
		return
	}
	c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
}

func SuperDelOrder(c *gin.Context) {
	u := new(ber)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	if u.Tp {
		go func() {
			var order []model.CoreQueryOrder
			model.DB().Where("`date` < ?", u.Date).Find(&order)

			tx := model.DB().Begin()
			for _, i := range order {
				model.DB().Where("work_id =?", i.WorkId).Delete(&model.CoreQueryOrder{})
				model.DB().Where("work_id =?", i.WorkId).Delete(&model.CoreQueryRecord{})
			}
			tx.Commit()
		}()
	} else {
		go func() {
			var order []model.CoreSqlOrder
			model.DB().Where("`date` < ?", u.Date).Find(&order)
			tx := model.DB().Begin()
			for _, i := range order {
				tx.Where("work_id =?", i.WorkId).Delete(&model.CoreSqlOrder{})
				tx.Where("work_id =?", i.WorkId).Delete(&model.CoreRollback{})
				tx.Where("work_id =?", i.WorkId).Delete(&model.CoreSqlRecord{})
			}
			tx.Commit()
		}()
	}
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_DELETE))
	return
}
