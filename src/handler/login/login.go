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

package login

import (
	"Yearning-go/internal/pkg/stringx"
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	SOURCE_LDAP = "LDAP"
)

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func userLdapLogin(u *loginForm, c *gin.Context) bool {
	attr, err := lib.LdapContent(&model.GloLdap, u.Username, u.Password)
	if err != nil {
		log.Println("LdapContent failed:", err.Error())
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(errors.New("远程登录失败")))
		return false
	}

	rule := "guest"
	if "" != model.GloLdap.Admin {
		conf := model.GloLdap
		conf.Type = model.GloLdap.Admin
		if _, err = lib.LdapContent(&conf, u.Username, u.Password); nil == err {
			rule = "admin"
		}
	}

	var account model.CoreAccount
	if model.DB().Where("username = ?", u.Username).First(&account).RecordNotFound() {
		model.DB().Create(&model.CoreAccount{
			Username:   u.Username,
			RealName:   stringx.Coalesce(attr.Name, "请重置你的真实姓名"),
			Password:   lib.DjangoEncrypt(lib.GenWorkid(), string(lib.GetRandom())),
			Rule:       rule,
			Department: "all",
			Email:      attr.Email,
			Mobile:     attr.Mobile,
			OpenId:     attr.OpenId,
			Source:     SOURCE_LDAP,
		})
		ix, _ := json.Marshal([]string{})
		model.DB().Create(&model.CoreGrained{Username: u.Username, Group: ix})
	}
	token, tokenErr := lib.JwtAuth(u.Username, account.Rule)
	if tokenErr != nil {
		// c.Logger().Error(tokenErr.Error())
		return false
	}
	dataStore := map[string]string{
		"token":       token,
		"permissions": account.Rule,
		"real_name":   account.RealName,
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(dataStore))
	return true
}

func UserGeneralLogin(c *gin.Context) {
	u := new(loginForm)
	if err := c.MustBindWith(u, lib.Binding{}); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}
	var account model.CoreAccount
	if model.DB().Where("username = ?", u.Username).First(&account).RecordNotFound() {
		userLdapLogin(u, c)
		return
	}

	if account.Username != u.Username {
		c.JSON(http.StatusOK, commom.ERR_LOGIN)
		return
	}
	if SOURCE_LDAP == account.Source {
		if !lib.LdapLogin(&model.GloLdap, u.Username, u.Password) {
			c.JSON(http.StatusOK, commom.ERR_LOGIN)
			return
		}
	} else if !lib.DjangoCheckPassword(&account, u.Password) {
		c.JSON(http.StatusOK, commom.ERR_LOGIN)
		return
	}

	token, tokenErr := lib.JwtAuth(u.Username, account.Rule)
	if tokenErr != nil {
		// c.Logger().Error(tokenErr.Error())
		return
	}
	dataStore := map[string]string{
		"token":       token,
		"permissions": account.Rule,
		"real_name":   account.RealName,
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(dataStore))
	return
}

func UserRegister(c *gin.Context) {

	if model.GloOther.Register {
		u := new(model.CoreAccount)
		if err := c.MustBindWith(u, lib.Binding{}); err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
			return
		}
		var unique model.CoreAccount
		ix, _ := json.Marshal([]string{})
		model.DB().Where("username = ?", u.Username).Select("username").First(&unique)
		if unique.Username != "" {
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(errors.New("用户已存在请重新注册！")))
			return
		}
		model.DB().Create(&model.CoreAccount{
			Username:   u.Username,
			RealName:   u.RealName,
			Password:   lib.DjangoEncrypt(u.Password, string(lib.GetRandom())),
			Rule:       "guest",
			Department: u.Department,
			Email:      u.Email,
		})
		model.DB().Create(&model.CoreGrained{Username: u.Username, Group: ix})
		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage("注册成功！"))
		return
	}
	c.JSON(http.StatusOK, commom.ERR_REGISTER)
}
