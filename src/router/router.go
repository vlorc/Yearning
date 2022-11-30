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

package router

import (
	"Yearning-go/src/apis"
	"Yearning-go/src/handler/login"
	autoTask2 "Yearning-go/src/handler/manager/autoTask"
	"Yearning-go/src/handler/manager/board"
	db2 "Yearning-go/src/handler/manager/db"
	flow2 "Yearning-go/src/handler/manager/flow"
	group2 "Yearning-go/src/handler/manager/group"
	proxy2 "Yearning-go/src/handler/manager/proxy"
	roles2 "Yearning-go/src/handler/manager/roles"
	"Yearning-go/src/handler/manager/settings"
	template2 "Yearning-go/src/handler/manager/template"
	user2 "Yearning-go/src/handler/manager/user"
	"Yearning-go/src/handler/openapi"
	audit2 "Yearning-go/src/handler/order/audit"
	"Yearning-go/src/handler/order/osc"
	query2 "Yearning-go/src/handler/order/query"
	"Yearning-go/src/handler/personal"
	"Yearning-go/src/lib"
	"Yearning-go/src/middleware"
	"Yearning-go/src/model"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

func SuperManageGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, role := lib.JwtParse(c)
		if role == "super" || focalPoint(c) {
			return
		}
		c.String(http.StatusForbidden, "非法越权操作！")
	}
}

func focalPoint(c *gin.Context) bool {

	if strings.Contains(c.Request.RequestURI, "/api/v2/manage/flow") && c.Request.Method == http.MethodPut {
		return true
	}

	if strings.Contains(c.Request.RequestURI, "/api/v2/manage/group") && c.Request.Method == http.MethodGet {
		return true
	}
	return false
}

func AuditGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, rule := lib.JwtParse(c)
		if rule != "guest" {
			return
		}
		c.String(http.StatusForbidden, "非法越权操作！")
	}
}

func AddRouter(e gin.IRouter, box fs.FS) {
	if os.Getenv("DEV") == "" {
		s, err := fs.ReadFile(box, "index.html")
		if err != nil {
			log.Fatal(err)
		}
		e.GET("/", func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8", s)
		})
	}
	e.POST("/login", login.UserGeneralLogin)
	e.POST("/register", login.UserRegister)
	e.GET("/fetch", login.UserReqSwitch)

	r := e.Group("/api/v2", middleware.JWTWithConfig(middleware.JwtConfig{SigningKey: []byte(model.JWT)}))
	personal.PersonalRestFulAPis().Route(r, "/common/:tp")
	apis.YearningDashApis().Route(r, "/dash/:tp")
	apis.YearningFetchApis().Route(r, "/fetch/:tp")
	apis.YearningQueryApis().Route(r, "/query/:tp")

	audit := r.Group("/audit", AuditGroup())
	audit2.AuditRestFulAPis().Route(audit, "/order/:tp")
	osc.AuditOSCFetchStateApis().Route(audit, "/osc/:work_id")
	query2.AuditQueryRestFulAPis().Route(audit, "/query/:tp")

	manager := r.Group("/manage", SuperManageGroup())
	manager.POST("/board/post", board.GeneralPostBoard)

	db := manager.Group("/db")
	db2.ManageDbApis().Route(db, "")

	account := manager.Group("/user")
	user2.SuperUserApis().Route(account, "")

	flow := manager.Group("/flow")
	flow2.FlowRestApis().Route(flow, "")

	template := manager.Group("/template")
	template2.TemplatesApis().Route(template, "")

	proxy := manager.Group("/proxy")
	proxy2.ProxyApis().Route(proxy, "")

	group := manager.Group("/group")
	group2.GroupsApis().Route(group, "")

	setting := manager.Group("/setting")
	settings.SettingsApis().Route(setting, "")

	roles := manager.Group("/roles")
	roles2.RolesApis().Route(roles, "")

	autoTask := manager.Group("/task")
	autoTask2.SuperAutoTaskApis().Route(autoTask, "")

	api := e.Group("/openapi/v1", middleware.NewOpenApi(model.C.Openapi.Appid, model.C.Openapi.Secret))
	api.POST("order/audit", openapi.AuditOrder)

}
