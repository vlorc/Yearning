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

package autoTask

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuperFetchAutoTaskList(c *gin.Context) {
	u := new(fetchAutoTask)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	var task []model.CoreAutoTask
	var pg int
	start, end := lib.Paging(u.Page, 15)
	if u.Find.Valve {
		model.DB().Model(model.CoreAutoTask{}).Scopes(commom.AccordingToOrderName(u.Find.Text)).Order("id desc").Count(&pg).Offset(start).Limit(end).Find(&task)
	} else {
		model.DB().Model(model.CoreAutoTask{}).Order("id desc").Count(&pg).Offset(start).Limit(end).Find(&task)
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(commom.CommonList{Data: task, Page: pg}))
}

func SuperDeleteAutoTask(c *gin.Context) {
	id := c.Query("id")
	model.DB().Where("id =?", id).Delete(&model.CoreAutoTask{})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_DELETE))
}

func SuperAutoTaskCreateOrEdit(c *gin.Context) {
	u := new(fetchAutoTask)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	switch u.Tp {
	case "create":
		u.Create()
	case "edit":
		u.Edit()
	case "active":
		u.Activation()
	}
	c.JSON(http.StatusOK,u.Resp)
}
