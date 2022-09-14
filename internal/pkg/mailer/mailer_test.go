package mailer

import (
	"testing"
)

func TestSend(t *testing.T) {
	err := Send(
		Config{
			Addr: "smtp.exmail.qq.com:465",
			User: "xxxxx@xxxxx.com",
			Pass: "xxxxx",
		},
		Message{
			To:      []string{"xxxxx@xxxxx.com"},
			ReplyTo: "",
			Subject: "Subject",
			Body:    "<p>hello</p>",
			Type:    "text/html",
		})

	if nil != err {
		t.Error(err.Error())
	}
}
