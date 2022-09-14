package ding

import (
	"Yearning-go/internal/pkg/messagex"
	"Yearning-go/internal/pkg/webhook"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

var typeMapping = map[messagex.Type]string{
	messagex.TYPE_TEXT:        TYPE_TEXT,
	messagex.TYPE_MARKDOWN:    TYPE_MARKDOWN,
	messagex.TYPE_LINK:        TYPE_LINK,
	messagex.TYPE_ACTION_CARD: TYPE_ACTION_CARD,
	messagex.TYPE_FEED_CARD:   TYPE_FEED_CARD,
}

func IsDingUrl(url string) bool {
	return strings.HasPrefix(url, "https://oapi.dingtalk.com")
}

func Send(conf Config, msg messagex.Message) error {
	c := webhook.Config{
		Url:    conf.Url,
		Method: "POST",
		Type:   "application/json; charset=utf-8",
	}
	if "" == conf.Url || "ding" == conf.Url {
		c.Url = "https://oapi.dingtalk.com/robot/send?access_token=" + conf.Token
	}
	if "" != conf.Secret {
		c.Url = sign(conf.Secret, c.Url)
	}
	msgtype := typeMapping[msg.Type]
	if "" == msgtype {
		msgtype = TYPE_TEXT
	}

	m := map[string]interface{}{
		"msgtype": msgtype,
		msgtype: message{
			Title:             msg.Subject,
			Text:              getText(msgtype, msg.Body),
			Content:           getContent(msgtype, msg.Body),
			Url:               msg.Link,
			ImageUrl:          msg.ImageUrl,
			ButtonOrientation: "",
			SingleTitle:       "",
			SingleURL:         "",
			Button:            nil,
			Links:             nil,
		},
	}

	if at := getAt(msg.Target); nil != at {
		m["at"] = at
	}

	b, _ := json.Marshal(m)

	return webhook.Send(c, bytes.NewReader(b))
}

func sign(secret, rawurl string) string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)

	return fmt.Sprintf("%s&timestamp=%d&sign=%s", rawurl, timestamp, sign)
}

func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getText(typ string, body string) string {
	if TYPE_TEXT != typ {
		return body
	}
	return ""
}

func getContent(typ string, body string) string {
	if TYPE_TEXT == typ {
		return body
	}
	return ""
}

func getAt(target messagex.Target) *At {
	if 0 == len(target.Mobiles) && 0 == len(target.OpenIds) && !target.All {
		return nil
	}
	if target.All {
		return &At{
			All: target.All,
		}
	}
	return &At{
		Mobiles: target.Mobiles,
		UserIds: target.OpenIds,
	}
}
