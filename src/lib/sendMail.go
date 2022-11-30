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

package lib

import (
	"Yearning-go/internal/pkg/ding"
	"Yearning-go/internal/pkg/mailer"
	"Yearning-go/internal/pkg/messagex"
	"Yearning-go/internal/pkg/qywx"
	"Yearning-go/internal/pkg/webhook"
	"Yearning-go/src/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"strings"
	"text/template"
	"time"
)

type UserInfo struct {
	ToUser  string
	User    string
	Pawd    string
	Smtp    string
	PubName string
}

type EventName string

const (
	EVENT_TEST               EventName = "SYSTEM_TEST"
	EVENT_ORDER_EXEC_CREATE  EventName = "ORDER_EXEC_CREATE"
	EVENT_ORDER_EXEC_PASS    EventName = "ORDER_EXEC_PASS"
	EVENT_ORDER_EXEC_REJECT  EventName = "ORDER_EXEC_REJECT"
	EVENT_ORDER_EXEC_SUCCESS EventName = "ORDER_EXEC_SUCCESS"
	EVENT_ORDER_EXEC_FAILED  EventName = "ORDER_EXEC_FAILED"
	EVENT_ORDER_EXEC_PERFORM EventName = "ORDER_EXEC_PERFORM"
	EVENT_ORDER_EXEC_UNDO    EventName = "ORDER_EXEC_UNDO"
	EVENT_ORDER_QUERY_CREATE EventName = "ORDER_QUERY_CREATE"
	EVENT_ORDER_QUERY_PASS   EventName = "ORDER_QUERY_PASS"
	EVENT_ORDER_QUERY_REJECT EventName = "ORDER_QUERY_REJECT"

	ORDER_EXEC_PREFIX    = "ORDER_EXEC_"
	ORDER_QUERY_PREFIX   = "ORDER_QUERY_"
	ORDER_PASS_SUFFIX    = "_PASS"
	ORDER_REJECT_suffix  = "_REJECT"
	ORDER_SUCCESS_SUFFIX = "_SUCCESS"
	ORDER_FAILED_suffix  = "_FAILED"
)

var eventNameMapping = map[EventName]string{
	EVENT_ORDER_EXEC_CREATE:  "工单提交",
	EVENT_ORDER_EXEC_PASS:    "工单通过",
	EVENT_ORDER_EXEC_REJECT:  "工单拒绝",
	EVENT_ORDER_EXEC_SUCCESS: "工单执行成功",
	EVENT_ORDER_EXEC_FAILED:  "工单执行失败",
	EVENT_ORDER_EXEC_PERFORM: "工单待审批",
	EVENT_ORDER_EXEC_UNDO:    "工单撤销",
	EVENT_ORDER_QUERY_CREATE: "查询申请创建",
	EVENT_ORDER_QUERY_PASS:   "查询申请通过",
	EVENT_ORDER_QUERY_REJECT: "查询申请拒绝",
}

var QueryTemplate = func(EventName) []model.CoreTemplate {
	return nil
}

func SendMail(rawurl string, msg messagex.Message) error {
	u, err := url.Parse(rawurl)
	if nil != err {
		return err
	}
	conf := mailer.Config{
		Addr:     u.Host,
		Ssl:      "smtps" == u.Scheme || u.Query().Get("ssl") == "true",
		Insecure: u.Query().Get("insecure") == "true",
		Timeout:  10,
	}
	if nil != u.User {
		conf.User = u.User.Username()
		conf.Pass, _ = u.User.Password()
	}

	return mailer.Send(conf, msg)
}

type TemplateParam struct {
	User     model.CoreAccount  `json:"user"`
	Order    TemplateOrderParam `json:"order"`
	Assigned model.CoreAccount  `json:"assigned"`
	Link     string             `json:"link,omitempty"`
	Kind     string             `json:"kind,omitempty"`
	Event    string             `json:"event,omitempty"`
	Name     string             `json:"name,omitempty"`
	Title    string             `json:"title,omitempty"`
	Date     string             `json:"date,omitempty"`
	Time     string             `json:"time,omitempty"`
	Host     string             `json:"host,omitempty"`
}

type TemplateOrderParam struct {
	ID          uint      `json:"id"`
	WorkId      string    `json:"work_id"`
	Username    string    `json:"username"`
	Status      uint      `json:"status"`
	Type        uint      `json:"type"` // 1 dml  0 ddl 3 query
	Backup      uint      `json:"backup"`
	IDC         string    `json:"idc"`
	Source      string    `json:"source"`
	DataBase    string    `json:"data_base"`
	Table       string    `json:"table"`
	Date        string    `json:"date"`
	SQL         string    `json:"sql"`
	Text        string    `json:"text"`
	Assigned    string    `json:"assigned"`
	Delay       string    `json:"delay"`
	RealName    string    `json:"real_name"`
	Executor    string    `json:"executor"`
	ExecuteTime string    `json:"execute_time"`
	Time        string    `json:"time"`
	IsKill      uint      `json:"is_kill"`
	CurrentStep int       `json:"current_step"`
	Percent     int       `json:"percent"`
	Current     int       `json:"current"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Export      uint      `json:"export"`
	QueryPer    int       `json:"query_per"`
	ExDate      string    `json:"ex_date"`
}

func MessagePush(workid string, event EventName, link string) {
	temp := QueryTemplate(event)
	if len(temp) == 0 {
		log.Println("can not match template by event: ", string(event))
		return
	}

	var param = TemplateParam{
		Link:  link,
		Event: string(event),
		Name:  eventNameMapping[event],
		Date:  time.Now().Format("2006-01-02"),
		Time:  time.Now().Format("15:04:05"),
		Host:  model.Host,
	}
	if strings.HasPrefix(string(event), ORDER_EXEC_PREFIX) {
		param.Kind = "提交人"
		var order model.CoreSqlOrder
		model.DB().Where("work_id =?", workid).First(&order)
		param.Order = TemplateOrderParam{
			ID:          order.ID,
			WorkId:      order.WorkId,
			Username:    order.Username,
			Status:      order.Status,
			Type:        order.Type,
			Backup:      order.Backup,
			IDC:         order.IDC,
			Source:      order.Source,
			DataBase:    order.DataBase,
			Table:       order.Table,
			Date:        order.Date,
			SQL:         order.SQL,
			Text:        order.Text,
			Assigned:    order.Assigned,
			Delay:       order.Delay,
			RealName:    order.RealName,
			Executor:    order.Executor,
			ExecuteTime: order.ExecuteTime,
			Time:        order.Time,
			IsKill:      order.IsKill,
			CurrentStep: order.CurrentStep,
			Percent:     order.Percent,
			Current:     order.Current,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}
	if strings.HasPrefix(string(event), ORDER_QUERY_PREFIX) {
		param.Kind = "提交人"
		var order model.CoreQueryOrder
		model.DB().Where("work_id =?", workid).First(&order)
		param.Order = TemplateOrderParam{
			ID:        order.ID,
			WorkId:    order.WorkId,
			Username:  order.Username,
			IDC:       order.IDC,
			Date:      order.Date,
			Text:      order.Text,
			Assigned:  order.Assigned,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}
	}
	if "" != param.Order.Username {
		var user model.CoreAccount
		model.DB().Where("username =?", param.Order.Username).First(&user)
		param.User = user
		param.User.Password = ""
	}
	if "" != param.Order.Assigned {
		var user model.CoreAccount
		model.DB().Where("username =?", param.Order.Assigned).First(&user)
		param.Assigned = user
		param.Assigned.Password = ""
	}

	for i := range temp {
		err := messagePushTemplate(&temp[i], &param)
		if nil != err {
			log.Println("push event template failed, event:", string(event), "template:", temp[i].Alias, "channel:", temp[i].Channel, "workId:", workid)
		} else {
			log.Println("push event template success, event:", string(event), "template:", temp[i].Alias, "channel:", temp[i].Channel, "workId:", workid)
		}
	}
}

func WebHookTestPush(conf model.Message) error {
	if "ding" == conf.WebHookUrl || ding.IsDingUrl(conf.WebHookUrl) {
		return ding.Send(
			ding.Config{
				Url:    conf.WebHookUrl,
				Token:  conf.Token,
				Secret: conf.Key,
			},
			messagex.Message{Body: "test"},
		)
	}
	if "qywx" == conf.WebHookUrl || qywx.IsQywxUrl(conf.WebHookUrl) {
		return qywx.Send(
			qywx.Config{
				Url:    conf.WebHookUrl,
				Token:  conf.Token,
				Secret: conf.Key,
			},
			messagex.Message{Body: "test"},
		)
	}
	return webhook.Send(webhook.Config{Url: conf.WebHookUrl}, strings.NewReader(`"event":"EVENT_TEST"`))
}

func messagePushTemplate(temp *model.CoreTemplate, param *TemplateParam) error {
	param.Title = temp.Title

	var b bytes.Buffer
	if "{{json}}" != temp.Body {
		t, err := template.New("").Parse(temp.Body)
		if nil != err {
			return errors.WithMessage(err, "Template parse")
		}
		if err = t.Execute(&b, param); nil != err {
			return errors.WithMessage(err, "Template execute")
		}
	} else {
		json.NewEncoder(&b).Encode(param)
	}

	msg := messagex.Message{
		Type:    messagex.Type(temp.Type),
		Subject: temp.Title,
		Body:    b.String(),
		Link:    param.Link,
	}
	if 0 != param.User.ID {
		msg.Sources = append(msg.Sources, messagex.Source{
			Name:   param.User.RealName,
			Email:  param.User.Email,
			OpenId: param.User.OpenId,
			Mobile: param.User.Mobile,
			Kind:   param.Kind,
		})
	}
	if strings.HasSuffix(param.Event, ORDER_PASS_SUFFIX) || strings.HasSuffix(param.Event, ORDER_REJECT_suffix) {
		appendTarget(&msg.Target, &param.User)
	} else if 0 != param.Assigned.ID {
		appendTarget(&msg.Target, &param.Assigned)
	}
	if strings.HasSuffix(param.Event, ORDER_SUCCESS_SUFFIX) || strings.HasSuffix(param.Event, ORDER_FAILED_suffix) {
		appendTarget(&msg.Target, &param.User)
	}

	if "email" == temp.Channel {
		return SendMail(temp.Url, msg)
	}
	if "ding" == temp.Channel || ding.IsDingUrl(temp.Url) {
		return ding.Send(ding.Config{Url: temp.Url, Secret: temp.Secret}, msg)
	}
	if "qywx" == temp.Channel || qywx.IsQywxUrl(temp.Url) {
		return qywx.Send(qywx.Config{Url: temp.Url}, msg)
	}
	if "webhook" == temp.Channel {
		return webhook.Send(webhook.Config{Url: temp.Url}, &b)
	}

	return fmt.Errorf("can not support channel %s", temp.Channel)
}

func appendTarget(target *messagex.Target, user *model.CoreAccount) {
	if 0 != user.ID {
		if "" != user.Email {
			target.Emails = append(target.Emails, user.Email)
		}
		if "" != user.Mobile {
			target.Mobiles = append(target.Mobiles, user.Mobile)
		}
		if "" != user.OpenId {
			target.OpenIds = append(target.OpenIds, user.OpenId)
		}
	}
}
