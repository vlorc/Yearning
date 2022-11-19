package proxy

func New(driver, rawurl, user, pass, secret string) Dialer {
	switch driver {
	case "ssh", "SSH":
		return NewSSH(rawurl, user, pass, secret)
	case "http", "https":
		return NewHttp(rawurl, user, pass, secret)
	case "httpx", "httpsx":
		return NewHttpx(rawurl, user, pass, secret)
	case "socks", "sockss":
		return NewSocks(rawurl, user, pass, secret)
	}
	return Direct{}
}
