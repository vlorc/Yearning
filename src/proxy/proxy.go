package proxy

import (
	"context"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"time"
)

type Dialer interface {
	proxy.ContextDialer
	proxy.Dialer
	Test(context.Context) error
	Driver() string
}

type Direct struct{}

type DialFunc func(context.Context, string, string) (net.Conn, error)

var _ proxy.ContextDialer = DialFunc(nil)
var _ proxy.Dialer = DialFunc(nil)

func (f DialFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return f(ctx, network, address)
}

func (f DialFunc) Dial(network, address string) (net.Conn, error) {
	return f(context.Background(), network, address)
}

func (Direct) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return proxy.Direct.DialContext(ctx, network, address)
}

func (f Direct) Dial(network, address string) (net.Conn, error) {
	return proxy.Direct.Dial(network, address)
}

func (f Direct) Test(ctx context.Context) error {
	return nil
}

func (f Direct) Driver() string {
	return "direct"
}

type ProxyServer struct {
	Name   string
	Target string
	Host   string
	Port   string
	Dialer Dialer
	ld     net.Listener
	begin  string
	count  uint32
}

func (s *ProxyServer) Close() error {
	if nil == s.ld {
		return nil
	}

	log.Println("Proxy close",
		"name:", s.Name,
		"driver:", s.Dialer.Driver(),
		"begin:", s.begin,
		"end:", time.Now().Format("2006-01-02 15:04:05"),
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
		"count:", s.count)

	time.AfterFunc(time.Minute*5, func() {
		_ = s.ld.Close()
		log.Println("Proxy destroy",
			"name:", s.Name,
			"driver:", s.Dialer.Driver(),
			"begin:", s.begin,
			"end:", time.Now().Format("2006-01-02 15:04:05"),
			"host:", s.Host,
			"port:", s.Port,
			"target:", s.Target,
			"count:", s.count,
		)
	})

	return nil
}

func (s *ProxyServer) Run() {
	if nil != s.ld {
		return
	}

	s.begin = time.Now().Format("2006-01-02 15:04:05")
	addr := net.JoinHostPort(s.Host, s.Port)
	l, err := net.Listen("tcp", addr)
	if nil != err {

		return
	}

	if "" == s.Port {
		_, port, _ := net.SplitHostPort(l.Addr().String())
		s.Port = port
	}

	s.ld = l

	go s.Serve()
}

func (s *ProxyServer) Serve() {
	if nil == s.ld {
		return
	}

	log.Println("Proxy listen begin",
		"name:", s.Name,
		"driver:", s.Dialer.Driver(),
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
	)

	for {
		conn, err := s.ld.Accept()
		if nil != err {
			log.Println("Proxy accept name:", s.Name, "driver:", s.Dialer.Driver(), err.Error())
			break
		}

		s.count++

		log.Println("Proxy accept",
			"name:", s.Name,
			"driver:", s.Dialer.Driver(),
			"begin:", s.begin,
			"host:", s.Host,
			"port:", s.Port,
			"target:", s.Target,
			"remote:", conn.RemoteAddr().String(),
		)

		go s.forward(conn)
	}

	log.Println("Proxy listen end",
		"name:", s.Name,
		"driver:", s.Dialer.Driver(),
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
	)
}

func (s *ProxyServer) forward(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	dst, err := s.Dialer.Dial("tcp", s.Target)
	if nil != err {
		_ = conn.Close()
		log.Println("Proxy dial failed",
			"name:", s.Name,
			"driver:", s.Dialer.Driver(),
			"begin:", s.begin,
			"host:", s.Host,
			"port:", s.Port,
			"target:", s.Target,
			"remote:", remote,
			"err:", err.Error(),
		)
		return
	}

	log.Println("Proxy forward begin",
		"name:", s.Name,
		"driver:", s.Dialer.Driver(),
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
		"remote:", remote,
	)
	go transfer(conn, dst)
	_ = transfer(dst, conn)

	log.Println("Proxy forward end",
		"name:", s.Name,
		"driver:", s.Dialer.Driver(),
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
		"remote:", remote,
	)
}

func transfer(dst io.WriteCloser, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}
