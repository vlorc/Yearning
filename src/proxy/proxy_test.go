package proxy

import (
	"golang.org/x/net/proxy"
	"testing"
	"time"
)

func TestProxySever(t *testing.T) {
	s := &ProxyServer{
		Name:   "test",
		Target: "127.0.0.1:80",
		Host:   "127.0.0.1",
		Dial:   proxy.Direct,
	}
	s.Run()

	time.Sleep(time.Second)

	t.Log("port:", s.Port)

	time.Sleep(time.Second * 15)

	_ = s.Close()
}
