package qywx

import (
	"Yearning-go/internal/pkg/messagex"
	"testing"
)

func TestSend(t *testing.T) {
	err := Send(
		Config{
			Token: "",
		},
		messagex.Message{
			Type: TYPE_TEXT,
			Body: "<font color=\"info\">绿色</font>\n<font color=\"comment\">灰色</font>\n<font color=\"warning\">橙红色</font>",
			Target: messagex.Target{
				Mobiles: []string{},
			},
		},
	)
	if nil != err {
		t.Error(err.Error())
	}
}
