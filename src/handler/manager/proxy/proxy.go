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
	"Yearning-go/src/proxy"
	"Yearning-go/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ProxyPage(c *gin.Context) {
	f := new(commom.PageInfo)
	if err := c.Bind(f); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	start, end := lib.Paging(f.Page, 10)
	page, data := service.ProxyService{}.Page(start, end)

	c.JSON(http.StatusOK, commom.SuccessPayload(
		commom.CommonList{
			Page: page,
			Data: data,
		},
	))
}

func ProxyUpdate(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	if user != "admin" {

		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}

	req := new(CommonProxyPost)
	if err := c.Bind(req); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	if "test" == req.Tp {
		dialer := &proxy.SSHDialer{
			Addr:     req.Proxy.Url,
			User:     req.Proxy.Username,
			Password: req.Proxy.Password,
			Secret:   req.Proxy.Secret,
		}
		info := service.ProxyService{}.InfoById(req.Proxy.ID)
		if nil != info && lib.IsHash(dialer.Password) {
			dialer.Password = info.Password
		}
		if nil != info && lib.IsHash(dialer.Secret) {
			dialer.Password = info.Secret
		}
		if err := dialer.Test(); nil == err {
			c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage("连接成功"))
		} else {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		}
		return
	}

	if 0 != req.Proxy.ID {
		service.ProxyService{}.Modify(&req.Proxy)
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(TEMPLATE_CREATE_SUCCESS, req.Proxy.Name)))
	} else {
		service.ProxyService{}.Create(&req.Proxy)
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(TEMPLATE_EDIT_SUCCESS, req.Proxy.Name)))
	}
}

func ProxyDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 63)
	if 0 == id {
		c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
	}
	info := service.ProxyService{}.DetailById(uint(id))
	if nil == info {
		c.JSON(http.StatusOK, commom.ERR_SOAR_ALTER_MESSAGE(TEMPLATE_NOT_EXIST))
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(info))
}
