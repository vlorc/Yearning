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

package personal

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/handler/fetch"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	ser "Yearning-go/src/parser"
	pb "Yearning-go/src/proto"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func ReferQueryOrder(c *gin.Context, user *string) {
	var u model.CoreAccount
	var t model.CoreQueryOrder

	d := new(commom.QueryOrder)
	if err := c.Bind(d); err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	state := 1

	if model.GloOther.Query {
		state = 2
	}

	model.DB().Select("real_name").Where("username =?", user).First(&u)

	if model.DB().Model(model.CoreQueryOrder{}).Where("username =? and query_per =?", user, 2).First(&t).RecordNotFound() {

		work := lib.GenWorkid()

		model.DB().Create(&model.CoreQueryOrder{
			WorkId:   work,
			Username: *user,
			Date:     time.Now().Format("2006-01-02 15:04"),
			Text:     d.Text,
			Assigned: d.Assigned,
			Export:   d.Export,
			IDC:      d.IDC,
			QueryPer: state,
			Realname: u.RealName,
			ExDate:   time.Now().Format("2006-01-02 15:04"),
		})

		if state == 2 {
			lib.MessagePush(work, lib.EVENT_ORDER_QUERY_CREATE, "")
		}

		c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_CREATE))
		return
	}
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_DUP))
}

func FetchQueryStatus(c *gin.Context) {

	user, _ := lib.JwtParse(c)

	var d model.CoreQueryOrder

	model.DB().Where("username =?", user).Last(&d)

	if lib.TimeDifference(d.ExDate) {
		model.DB().Model(model.CoreQueryOrder{}).Where("username =?", user).Update(&model.CoreQueryOrder{QueryPer: 3})
	}
	c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"status": d.QueryPer, "export": model.GloOther.Export, "idc": d.IDC}))
	return
}

func FetchQueryDatabaseInfo(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	var d model.CoreQueryOrder
	var u model.CoreDataSource

	model.DB().Where("username =?", user).Last(&d)

	if d.QueryPer == 1 {

		source := new(commom.QueryOrder)

		if err := c.Bind(source); err != nil {
			//c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
			return
		}

		model.DB().Where("source =?", source.Source).First(&u)

		result, err := commom.ScanDataRows(u, "", "SHOW DATABASES;", "库名", true)

		if err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
			return
		}

		var info []map[string]interface{}

		info = append(info, map[string]interface{}{
			"title":    source.Source,
			"expand":   "true",
			"children": result.BaseList,
		})
		c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"info": info, "status": d.Export, "highlight": result.Highlight, "sign": fetch.FetchTplAuditor(u.IDC), "idc": u.IDC}))
		return
	} else {
		c.JSON(http.StatusOK, commom.SuccessPayload(0))
	}
}

func FetchQueryTableInfo(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	t := c.Query("title")
	// todo source改方法 不然中文无法识别
	source := c.Query("source")
	unescape, _ := url.QueryUnescape(source)
	var d model.CoreQueryOrder
	var u model.CoreDataSource
	model.DB().Where("username =?", user).Last(&d)

	if d.QueryPer == 1 {

		model.DB().Where("source =?", unescape).First(&u)

		result, err := commom.ScanDataRows(u, t, "SHOW TABLES;", "表名", true)
		if err != nil {
			// c.Logger().Error(err.Error())
			c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
			return
		}
		c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"table": result.Query, "highlight": result.Highlight}))
		return

	} else {
		c.JSON(http.StatusOK, commom.SuccessPayload(0))
	}
}

func FetchQueryTableStruct(c *gin.Context) {
	t := new(queryBind)
	if err := c.Bind(t); err != nil {
		c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		return
	}

	unescape, _ := url.QueryUnescape(t.Source)
	user, _ := lib.JwtParse(c)
	var d model.CoreQueryOrder
	var u model.CoreDataSource
	var f []ser.FieldInfo
	model.DB().Where("username =?", user).Last(&d)
	model.DB().Where("source =?", unescape).First(&u)
	ps := lib.Decrypt(u.Password)
	host := u.IP
	if "" != u.Proxy {
		host = strings.Join([]string{lib.PROXY_PREFIX, u.Proxy, "_", hex.EncodeToString([]byte(net.JoinHostPort(u.IP, strconv.Itoa(u.Port))))}, "")
	}

	db, e := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", u.Username, ps, host, u.Port, t.DataBase))
	if e != nil {
		// c.Logger().Error(e.Error())
		c.JSON(http.StatusInternalServerError, commom.SuccessPayLoadToMessage(ER_DB_CONNENT))
		return
	}
	defer db.Close()

	if err := db.Raw(fmt.Sprintf("SHOW FULL FIELDS FROM `%s`.`%s`", t.DataBase, t.Table)).Scan(&f).Error; err != nil {
		// c.Logger().Error(err.Error())
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		return
	}

	c.JSON(http.StatusOK, commom.SuccessPayload(f))
}

func FetchQueryResults(c *gin.Context, user *string) {

	req := new(lib.QueryDeal)

	clock := time.Now()

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	//需自行实现查询SQL LIMIT限制
	err := lib.Limit(req, &pb.LibraAuditOrder{SQL: req.Sql}, "")

	if err != nil {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		return
	}

	var d model.CoreQueryOrder

	model.DB().Where("username =? AND query_per =?", user, 1).Last(&d)

	if lib.TimeDifference(d.ExDate) {
		model.DB().Model(model.CoreQueryOrder{}).Where("username =?", user).Update(&model.CoreQueryOrder{QueryPer: 3})
		c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"status": true}))
		return
	}

	//结束
	data := new(lib.Query)

	var u model.CoreDataSource

	model.DB().Where("source =?", req.Source).First(&u)

	err = data.QueryRun(&u, req)

	if err != nil {
		c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		return
	}

	queryTime := int(time.Since(clock).Seconds() * 1000)

	go func(w, s string, ex int) {
		model.DB().Create(&model.CoreQueryRecord{SQL: s, WorkId: w, ExTime: ex, Time: time.Now().Format("2006-01-02 15:04"), Source: req.Source, BaseName: req.DataBase})
	}(d.WorkId, req.Sql, queryTime)

	c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"title": data.Field, "data": data.Data, "status": false, "time": queryTime, "total": len(data.Data)}))
	return
}

func UndoQueryOrder(c *gin.Context) {
	user, _ := lib.JwtParse(c)
	model.DB().Model(model.CoreQueryOrder{}).Where("username =?", user).Update(map[string]interface{}{"query_per": 3})
	c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(commom.ORDER_IS_END))
}
