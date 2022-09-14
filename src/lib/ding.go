package lib

import (
	"Yearning-go/internal/pkg/ding"
	"Yearning-go/internal/pkg/messagex"
	"Yearning-go/src/model"
	"log"
)

func SendDing(msg model.Message, sv string) {
	err := ding.Send(
		ding.Config{
			Url:    msg.WebHookUrl,
			Secret: msg.Key,
		},
		messagex.Message{
			Type:    messagex.TYPE_MARKDOWN,
			Subject: "Yearning sql审计平台",
			Body:    sv,
		},
	)

	if err != nil {
		log.Println(err.Error())
	}
}
