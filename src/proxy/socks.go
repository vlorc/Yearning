package proxy

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

type SocksDialer struct {
	dialer    proxy.Dialer
	ctxDialer proxy.ContextDialer
}

func NewSocks(rawurl, user, pass, secret string) Dialer {
	var auth *proxy.Auth
	if "" != user {
		auth = &proxy.Auth{
			User:     user,
			Password: pass,
		}
	}

	var parent proxy.Dialer = proxy.Direct

	if conf := secretToTlsConfig(secret); nil != conf {
		parent = DialFunc(func(ctx context.Context, s string, s2 string) (net.Conn, error) {
			return tls.Dial("tcp", rawurl, conf)
		})
	}
	dialer, err := proxy.SOCKS5("tcp", rawurl, auth, parent)
	if nil != err {
		log.Println("proxy socks create failed:", err.Error())
		return nil
	}

	ctxDialer, _ := dialer.(proxy.ContextDialer)

	return &SocksDialer{
		dialer:    dialer,
		ctxDialer: ctxDialer,
	}
}

func (s *SocksDialer) Dial(network, addr string) (net.Conn, error) {
	return s.dialer.Dial(network, addr)
}

func (s *SocksDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if nil != s.ctxDialer {
		return s.ctxDialer.DialContext(ctx, network, addr)
	}
	return s.dialer.Dial(network, addr)
}

func (s *SocksDialer) Test(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
