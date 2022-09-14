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

package template

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TemplatePage(c *gin.Context) {
	f := new(commom.PageInfo)
	if err := c.Bind(f); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	start, end := lib.Paging(f.Page, 10)
	page, data := service.TemplateService{}.Page(start, end)

	c.JSON(http.StatusOK, commom.SuccessPayload(
		commom.CommonList{
			Page:    page,
			Data:    data,
		},
	))
}

func TemplateUpdate(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	if user == "admin" {
		req := new(CommonTemplatePost)
		if err := c.Bind(req); err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
			return
		}
		if 0 != req.Template.ID {
			service.TemplateService{}.Modify(&req.Template)
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(TEMPLATE_CREATE_SUCCESS, req.Template.Name)))
		} else {
			service.TemplateService{}.Create(&req.Template)
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(TEMPLATE_EDIT_SUCCESS, req.Template.Name)))
		}
	}
	c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
}

func TemplateDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 63)
	if 0 == id {
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
	info := service.TemplateService{}.InfoById(uint(id))
	if nil == info {
		c.JSON(http.StatusOK, commom.ERR_SOAR_ALTER_MESSAGE(TEMPLATE_NOT_EXIST))
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(info))
}
