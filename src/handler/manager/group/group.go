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

package group

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func SuperGroup(c *gin.Context) {
	var page int
	var roles []model.CoreRoleGroup

	f := new(commom.PageInfo)
	if err := c.Bind(f); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	start, end := lib.Paging(f.Page, 10)
	var source []model.CoreDataSource
	var query []model.CoreDataSource
	var u []model.CoreAccount
	model.DB().Select("source").Scopes(commom.AccordingToGroupSourceIsQuery(0, 2)).Find(&source)
	model.DB().Select("source").Scopes(commom.AccordingToGroupSourceIsQuery(1, 2)).Find(&query)
	model.DB().Select("username").Scopes(commom.AccordingToRuleSuperOrAdmin()).Find(&u)
	if f.Find.Valve {
		model.DB().Model(model.CoreRoleGroup{}).Scopes(commom.AccordingToOrderName(f.Find.Text)).Count(&page).Offset(start).Limit(end).Find(&roles)
	} else {
		model.DB().Model(model.CoreRoleGroup{}).Count(&page).Offset(start).Limit(end).Find(&roles)
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(
		commom.CommonList{
			Page:    page,
			Data:    roles,
			Source:  source,
			Query:   query,
			Auditor: u,
		},
	))
}

func SuperGroupUpdate(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	if user == "admin" {
		req := new(updateReq)
		if err := c.MustBindWith(req, lib.Binding{}); err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
			return
		}
		g, err := json.Marshal(req.Permission)
		if err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
			return
		}
		if req.Tp == 1 {
			var s model.CoreRoleGroup
			if model.DB().Scopes(commom.AccordingToNameEqual(req.Username)).First(&s).RecordNotFound() {
				model.DB().Create(&model.CoreRoleGroup{
					Name:        req.Username,
					Permissions: g,
				})
			} else {
				model.DB().Model(model.CoreRoleGroup{}).Scopes(commom.AccordingToNameEqual(req.Username)).Update(&model.CoreRoleGroup{Permissions: g})
			}
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(GROUP_CREATE_SUCCESS, req.Username)))
			return
		} else {
			g, _ := json.Marshal(req.Group)
			model.DB().Model(model.CoreGrained{}).Scopes(commom.AccordingToUsernameEqual(req.Username)).Updates(model.CoreGrained{Group: g})
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(GROUP_EDIT_SUCCESS, req.Username)))
			return
		}
	}
	c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
}

func SuperClearUserRule(c *gin.Context) {
	args := c.Query("clear")
	scape, _ := url.QueryUnescape(args)
	var j []model.CoreGrained
	var m1 []string
	model.DB().Scopes(commom.AccordingToGroupNameIsLike(scape)).Find(&j)
	for _, i := range j {
		_ = json.Unmarshal(i.Group, &m1)
		marshalGroup, _ := json.Marshal(lib.ResearchDel(m1, scape))
		model.DB().Model(model.CoreGrained{}).Scopes(commom.AccordingToUsernameEqual(i.Username)).Update(&model.CoreGrained{Group: marshalGroup})
	}
	model.DB().Model(model.CoreRoleGroup{}).Scopes(commom.AccordingToNameEqual(scape)).Delete(&model.CoreRoleGroup{})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(GROUP_DELETE_SUCCESS, scape)))
}

func SuperUserRuleMarge(c *gin.Context) {
	req := new(margeReq)
	if err := c.MustBindWith(req, lib.Binding{}); err != nil {
		//c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	m3 := lib.MultiUserRuleMarge(strings.Split(req.Group, ","))
	c.JSON(http.StatusOK, commom.SuccessPayload(m3))
	return
}
