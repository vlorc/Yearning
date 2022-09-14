package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gopkg.in/ldap.v3"
	"strings"
)

type Config struct {
	Url    string `json:"url"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Filter string `json:"filter"`
	Dn     string `json:"dn"`
	Ssl    bool   `json:"ssl"`
}

type Attr struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	OpenId string `json:"openId"`
}

type Request struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Test bool   `json:"test"`
}

func Connect(conf Config, req Request, attr *Attr) (err error) {
	var ld *ldap.Conn

	if conf.Ssl {
		ld, err = ldap.DialTLS("tcp", conf.Url, &tls.Config{InsecureSkipVerify: true})
	} else {
		ld, err = ldap.Dial("tcp", conf.Url)
	}
	if err != nil {
		return err
	}
	defer ld.Close()

	if "" != conf.User {
		if err = ld.Bind(conf.User, conf.Pass); err != nil || req.Test {
			return err
		}
	}
	if strings.Index(conf.Dn, "%s") > 0 {
		dn := fmt.Sprintf(conf.Dn, req.User)
		if err = ld.Bind(dn, req.Pass); err != nil || req.Test {
			return err
		}
		conf.Dn = dn
	}
	if "" == conf.Filter {
		conf.Filter = "(&(objectClass=inetOrgPerson)(uid=%s))"
	}

	var attrs []string
	if nil == attr {
		attrs = []string{"dn"}
	}

	sr, err := ld.Search(ldap.NewSearchRequest(
		conf.Dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(conf.Filter, req.User),
		attrs,
		nil,
	))
	if err != nil {
		return err
	}
	if len(sr.Entries) != 1 {
		return errors.New("User does not exist or too many entries returned")
	}

	if nil != attr {
		attr.Name = ldapAttr(sr.Entries[0], "nickname", "displayName", "name", "cn")
		attr.Mobile = ldapAttr(sr.Entries[0], "mobile", "homePhone", "phone")
		attr.Id = ldapAttr(sr.Entries[0], "uid", "uidNumber")
		attr.Avatar = ldapAttr(sr.Entries[0], "avatar", "jpegPhoto", "photo", "picture")
		attr.Email = ldapAttr(sr.Entries[0], "email", "mail")
		attr.OpenId = ldapAttr(sr.Entries[0], "openid")
	}

	return ld.Bind(sr.Entries[0].DN, req.Pass)
}

func ldapAttr(entry *ldap.Entry, keys ...string) string {
	for _, k := range keys {
		if v := entry.GetAttributeValue(k); "" != v {
			return v
		}
	}
	return ""
}
