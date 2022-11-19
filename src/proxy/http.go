package proxy

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type HttpDialer struct {
	conf    *tls.Config
	dialer  proxy.ContextDialer
	addr    string
	data    []byte
	isHttps bool
	code    int
}

const (
	textClientUpgradeRequest = "GET %s HTTP/1.1\nHost: %s\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Protocol: bitch\nSec-WebSocket-Version: 13\n%s%s%s\n"
	textClientConnectRequest = "CONNECT %s HTTP/1.1\nHost: %s\n%s%s%s\n"
)

func NewHttpx(rawurl, user, pass, secret string) Dialer {
	d := createHttp(rawurl, user, pass, secret, textClientUpgradeRequest)
	if nil != d {
		d.code = http.StatusSwitchingProtocols
	}
	return d
}

func NewHttp(rawurl, user, pass, secret string) Dialer {
	return createHttp(rawurl, user, pass, secret, textClientConnectRequest)
}

func createHttp(rawurl, user, pass, secret, format string) *HttpDialer {
	u, err := url.Parse(rawurl)
	if nil != err {
		log.Println("Proxy url parse failed",
			"url:", rawurl,
			"err:", err.Error(),
		)
		return nil
	}
	path := u.Path
	if "" == path {
		path = "/"
	}

	var data string
	if "" != user {
		data = fmt.Sprintf(format, path, u.Host, "Authorization: Basic ", base64.StdEncoding.EncodeToString([]byte(user+":"+pass)), "\n")
	} else {
		data = fmt.Sprintf(format, path, u.Host, "", "", "")
	}

	addr := u.Host
	if _, _, err = net.SplitHostPort(u.Host); nil != err {
		if "http" == u.Scheme {
			addr = u.Host + ":80"
		} else {
			addr = u.Host + ":443"
		}
	}

	d := &HttpDialer{
		conf:    secretToTlsConfig(secret),
		addr:    addr,
		isHttps: "https" == u.Scheme,
		data:    []byte(data),
		code:    http.StatusOK,
		dialer:  new(net.Dialer),
	}
	if d.isHttps {
		d.dialer = &tls.Dialer{
			NetDialer: new(net.Dialer),
			Config:    d.conf,
		}
	}

	return d
}

func (s *HttpDialer) Conn(ctx context.Context) (net.Conn, error) {
	conn, err := s.dialer.DialContext(ctx, "tcp", s.addr)
	if err != nil {
		err = errors.WithMessage(err, "proxy dial tcp")
		return nil, err
	}

	if _, err = conn.Write(s.data); err != nil {
		_ = conn.Close()
		err = errors.WithMessage(err, "proxy dial write request")
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if nil != err {
		_ = conn.Close()
		err = errors.WithMessage(err, "proxy dial read response")
		return nil, err
	}
	if resp.StatusCode != s.code {
		_ = conn.Close()
		err = errors.WithMessagef(err, "proxy dial response failed %d", resp.StatusCode)
		return nil, err
	}

	if !s.isHttps && nil != s.conf {
		log.Println("proxy http connect upgrade ready")
		tc := tls.Client(conn, s.conf)
		if err = tc.HandshakeContext(context.Background()); nil != err {
			_ = conn.Close()
			log.Println("proxy http connect upgrade error:", err.Error())
			return nil, err
		}
		log.Println("proxy http connect upgrade finish")
		return tc, nil
	}

	return conn, nil
}

func (s *HttpDialer) Test(ctx context.Context) error {
	conn, err := s.Conn(ctx)
	if nil != err {
		return err
	}
	return conn.Close()
}

func (s *HttpDialer) DialContext(ctx context.Context, _, _ string) (net.Conn, error) {
	return s.Conn(ctx)
}

func (s *HttpDialer) Dial(_, _ string) (conn net.Conn, err error) {
	return s.Conn(context.Background())
}

func secretToTlsConfig(secret string) *tls.Config {
	if "" == secret {
		return nil
	}

	buf, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Println("proxy secret decode error:", err.Error())
		return nil
	}

	z, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		log.Println("proxy http secret read error:", err.Error())
		return nil
	}

	var pem, key, ca []byte
	for i := range z.File {
		r, err := z.File[i].Open()
		if nil != err {
			continue
		}
		b, err := io.ReadAll(r)
		if nil != err {
			continue
		}
		switch z.File[i].Name {
		case "client.pem", "client.crt":
			pem = b
		case "client.key":
			key = b
		case "ca.pem", "ca.crt":
			ca = b
		}
	}
	conf := &tls.Config{
		RootCAs:            x509.NewCertPool(),
		InsecureSkipVerify: true,
		Certificates:       make([]tls.Certificate, 1),
	}

	if conf.Certificates[0], err = tls.X509KeyPair(pem, key); nil != err {
		log.Println("proxy http secret read x509 failed: ", err.Error())
		return nil
	}
	if len(ca) > 0 {
		conf.RootCAs.AppendCertsFromPEM(ca)
	}
	return conf
}
