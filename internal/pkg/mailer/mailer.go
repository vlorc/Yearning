package mailer

import (
	"Yearning-go/internal/pkg/messagex"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"time"
)

type Config struct {
	Addr     string `json:"url"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Ssl      bool   `json:"ssl"`
	Insecure bool   `json:"insecure"`
	Timeout  int    `json:"timeout"`
}

func send(c Config, m Message) (err error) {
	host, port, _ := net.SplitHostPort(c.Addr)

	var conn net.Conn

	dialer := &net.Dialer{Timeout: time.Second * time.Duration(c.Timeout)}
	if port == "465" || c.Ssl {
		conn, err = tls.DialWithDialer(dialer, "tcp", c.Addr, &tls.Config{InsecureSkipVerify: c.Insecure})
	} else {
		conn, err = dialer.Dial("tcp", c.Addr)
	}
	if nil != err {
		return err
	}

	cli, err := smtp.NewClient(conn, host)
	if nil != err {
		return err
	}
	defer cli.Close()

	if ok, _ := cli.Extension("STARTTLS"); ok {
		if err = cli.StartTLS(&tls.Config{InsecureSkipVerify: c.Insecure}); err != nil {
			return err
		}
	}

	if ok, _ := cli.Extension("AUTH"); ok {
		if err = cli.Auth(smtp.PlainAuth("", c.User, c.Pass, host)); nil != err {
			return err
		}
	}

	if err = cli.Mail(c.User); nil != err {
		return err
	}

	for _, addr := range m.Tolist() {
		if err = cli.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := cli.Data()
	if err != nil {
		return err
	}
	if "" == m.From.Address {
		m.From = mail.Address{Name: c.Name, Address: c.User}
	}
	if "" == m.Type {
		m.Type = TEXT
	}
	if _, err = io.Copy(w, bytes.NewReader(m.Bytes())); err != nil {
		return err
	}
	return w.Close()
}

func Send(c Config, m messagex.Message) (err error) {
	msg := Message{
		To:      m.Target.Emails,
		Subject: m.Subject,
		Body:    m.Body,
	}
	if messagex.TYPE_HTML == m.Type {
		msg.Type = HTML
	} else {
		msg.Type = TEXT
	}
	if len(m.Files) > 0 {
		msg.Attachments = make(map[string]*Attachment)
		for i := range m.Files {
			msg.Attachments[m.Files[i].Name] = &Attachment{
				Filename: m.Files[i].Name,
				Data:     m.Files[i].Data,
				Inline:   m.Files[i].Inline,
			}
		}
	}

	return send(c, msg)
}
