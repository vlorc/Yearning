package ldap

import "testing"

func TestLdap(t *testing.T) {
	attr := &Attr{}
	err := Connect(
		Config{
			Url:  "127.0.0.1:3890",
			User: "cn=admin,dc=example,dc=org",
			Pass: "admin",
			Dn:   "dc=example,dc=org",
			Ssl:  false,
		},
		Request{
			User: "qqq",
			Pass: "qqqbbb",
		},
		attr,
	)
	if nil != err {
		t.Errorf(err.Error())
	}
	t.Log(attr)
}
