package proxy

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

type SSHDialer struct {
	Addr     string
	User     string
	Password string
	Secret   string
}

type SSHConn struct {
	net.Conn
	client *ssh.Client
}

var _ proxy.Dialer = &SSHDialer{}

func NewSSH(rawurl, user, pass, secret string) Dialer {
	return &SSHDialer{
		Addr:     rawurl,
		User:     user,
		Password: pass,
		Secret:   secret,
	}
}

func (s *SSHDialer) auth() (result []ssh.AuthMethod) {
	result = make([]ssh.AuthMethod, 0)
	if "" == s.Secret {
		result = append(result, ssh.Password(s.Password))
		return
	}

	var signer ssh.Signer
	var err error
	if "" == s.Password {
		signer, err = ssh.ParsePrivateKey([]byte(s.Secret))
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(s.Secret), []byte(s.Password))
	}
	if err != nil {
		log.Println("SSH parse failed:", err.Error())
	} else {
		result = append(result, ssh.PublicKeys(signer))
	}

	return
}

func (s *SSHDialer) Test(ctx context.Context) error {
	client, err := s.Conn(ctx)
	if nil != err {
		return err
	}
	return client.Close()
}

func (s *SSHDialer) Conn(ctx context.Context) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: s.auth(),
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", s.Addr, config)
}

func (s *SSHDialer) Dial(network, addr string) (net.Conn, error) {
	return s.DialContext(context.Background(), network, addr)
}

func (s *SSHDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	client, err := s.Conn(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := client.Dial(network, addr)
	if nil != err {
		_ = client.Close()
		return nil, errors.WithMessage(err, "SSH dial")
	}

	return &SSHConn{Conn: conn, client: client}, nil
}

func (c *SSHConn) Close() error {
	err := c.Conn.Close()
	if nil != c.client {
		return c.client.Close()
	}
	return err
}

func (s *SSHDialer) Driver() string {
	return "ssh"
}
