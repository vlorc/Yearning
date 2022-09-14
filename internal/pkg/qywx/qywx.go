package qywx

import (
	"Yearning-go/internal/pkg/messagex"
	"Yearning-go/internal/pkg/webhook"
	"bytes"
	"encoding/json"
	"strings"
)

type Config struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

var typeMapping = map[messagex.Type]string{
	messagex.TYPE_TEXT:        TYPE_TEXT,
	messagex.TYPE_MARKDOWN:    TYPE_MARKDOWN,
	messagex.TYPE_IMAGE:       TYPE_IMAGE,
	messagex.TYPE_LINK:        TYPE_NEWS,
	messagex.TYPE_FILE:        TYPE_FILE,
	messagex.TYPE_ACTION_CARD: TYPE_TEMPLATE_CARD,
}

func IsQywxUrl(url string) bool {
	return strings.HasPrefix(url, "https://qyapi.weixin.qq.com")
}

func Send(conf Config, msg messagex.Message) error {
	c := webhook.Config{
		Url:    conf.Url,
		Method: "POST",
		Type:   "application/json; charset=utf-8",
	}
	if "" == conf.Url || "qywx" == conf.Url {
		c.Url = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + conf.Token
	}
	msgtype := typeMapping[msg.Type]
	if "" == msgtype {
		msgtype = TYPE_TEXT
		msg.Type = messagex.TYPE_TEXT
	}

	m := message{
		Type:         msgtype,
		Text:         getMessageText(msgtype, &msg),
		Markdown:     getMessageMarkdown(msgtype, &msg),
		TemplateCard: getMessageTemplateCard(msgtype, &msg),
	}

	b, _ := json.Marshal(m)

	return webhook.Send(c, bytes.NewReader(b))
}

func getText(msg *messagex.Message) *messageText {
	m := &messageText{
		Content: msg.Body,
		UserIds: msg.Target.OpenIds,
		Mobiles: msg.Target.Mobiles,
	}
	if msg.Target.All {
		m.UserIds = []string{"@all"}
		m.Mobiles = nil
	}
	return m
}

func getMessageText(msgtype string, msg *messagex.Message) *messageText {
	if TYPE_TEXT != msgtype {
		return nil
	}
	return getText(msg)
}

func getMessageMarkdown(msgtype string, msg *messagex.Message) *messageText {
	if TYPE_MARKDOWN != msgtype {
		return nil
	}
	return getText(msg)
}

func getMessageImage(msgtype string, msg *messagex.Message) *messageImage {
	if TYPE_IMAGE != msgtype {
		return nil
	}

	return &messageImage{
		Base64: "",
		Md5:    "",
	}
}

func getMessageTemplateCard(msgtype string, msg *messagex.Message) *messageTemplateCard {
	if TYPE_TEMPLATE_CARD != msgtype {
		return nil
	}

	m := &messageTemplateCard{
		CardType: "text_notice",
	}
	if nil != msg.App {
		m.Source = &messageTemplateCardSource{
			IconURL:   msg.App.IconUrl,
			Desc:      msg.App.Name,
			DescColor: 1,
		}
	}
	if "" != msg.Subject {
		m.MainTitle = &messageTemplateCardContent{
			Title: msg.Subject,
		}
	}
	if len(msg.Sources) > 0 {
		for i := range msg.Sources {
			if "" != msg.Sources[i].Kind && "" != msg.Sources[i].Name {
				m.HorizontalContentList = append(m.HorizontalContentList, messageTemplateCardHorizontalContent{
					Keyname: msg.Sources[i].Kind,
					Value:   msg.Sources[i].Name,
				})
			}
			if "" != msg.Sources[i].Title && "" != msg.Sources[i].Phrase {
				m.HorizontalContentList = append(m.HorizontalContentList, messageTemplateCardHorizontalContent{
					Keyname: msg.Sources[i].Title,
					Value:   msg.Sources[i].Phrase,
				})
			}
		}
	}

	return m
}
