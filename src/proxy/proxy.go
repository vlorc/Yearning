package proxy

import (
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"time"
)

type ProxyServer struct {
	Name   string
	Target string
	Host   string
	Port   string
	Dial   proxy.Dialer
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
		"begin:", s.begin,
		"end:", time.Now().Format("2006-01-02 15:04:05"),
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
		"count:", s.count)

	time.AfterFunc(time.Minute * 3, func() {
		_ = s.ld.Close()
		log.Println("Proxy destroy",
			"name:", s.Name,
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
		log.Println("Proxy listen failed",
			"name:", s.Name,
			"begin:", s.begin,
			"host:", s.Host,
			"port:", s.Port,
			"target:", s.Target,
			"err:", err.Error(),
		)
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
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
	)

	for {
		conn, err := s.ld.Accept()
		if nil != err {
			log.Println("Proxy accept name:", s.Name, err.Error())
			break
		}

		s.count++

		log.Println("Proxy accept",
			"name:", s.Name,
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
		"begin:", s.begin,
		"host:", s.Host,
		"port:", s.Port,
		"target:", s.Target,
	)
}

func (s *ProxyServer) forward(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	dst, err := s.Dial.Dial("tcp", s.Target)
	if nil != err {
		_ = conn.Close()
		log.Println("Proxy dial failed",
			"name:", s.Name,
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
