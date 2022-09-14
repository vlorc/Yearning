package lib

import (
	"Yearning-go/internal/pkg/ldap"
	"Yearning-go/src/model"
	"log"
)

func LdapLogin(l *model.Ldap, user string, pass string) bool {
	_, err := LdapContent(l, user, pass)
	return nil == err
}

func LdapContent(l *model.Ldap, user string, pass string) (*ldap.Attr, error) {
	attr := &ldap.Attr{}
	err := ldap.Connect(
		ldap.Config{
			Url:    l.Url,
			User:   l.User,
			Pass:   l.Password,
			Filter: l.Type,
			Dn:     l.Sc,
			Ssl:    l.Ldaps,
		},
		ldap.Request{
			User: user,
			Pass: pass,
		},
		attr,
	)

	if nil != err {
		log.Println("Ldap connect failed:", err.Error())
		return nil, err
	}

	return attr, nil
}

func LdapTest(l *model.Ldap) error {
	err := ldap.Connect(
		ldap.Config{
			Url:    l.Url,
			User:   l.User,
			Pass:   l.Password,
			Filter: l.Type,
			Dn:     l.Sc,
			Ssl:    l.Ldaps,
		},
		ldap.Request{Test: true},
		nil,
	)

	if nil != err {
		log.Println("Ldap test failed:", err.Error())
	}

	return err
}
