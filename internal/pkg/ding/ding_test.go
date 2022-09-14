package ding

import (
	"Yearning-go/internal/pkg/messagex"
	"testing"
)

func TestSend(t *testing.T) {
	err := Send(
		Config{
			Secret: "",
			Token:  "",
		},
		messagex.Message{
			Type:     "",
			Subject:  "",
			Body:     "test123",
			ImageUrl: "",
		},
	)
	if nil != err {
		t.Error(err.Error())
	}
}
