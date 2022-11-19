package service

import (
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"Yearning-go/src/proxy"
	"encoding/hex"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"net"
	"strings"
)

type ProxyService struct{}

func init() {
	mysql.RegisterDial("tcp", Dial)
}

func Dial(addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if nil != err {
		return nil, err
	}

	if !strings.HasPrefix(host, lib.PROXY_PREFIX) {
		return net.Dial("tcp", addr)
	}
	p := strings.LastIndex(host, "_")

	target, err := hex.DecodeString(host[p+1:])
	if nil != err {
		return nil, errors.WithMessage(err, "Host decode target")
	}
	alias := host[len(lib.PROXY_PREFIX):p]

	info := ProxyService{}.InfoByAlias(alias)
	if nil == info {
		return nil, fmt.Errorf("can not found proxy with alias '%s'", alias)
	}
	if 1 != info.Status {
		return nil, fmt.Errorf("proxy is disable with alias '%s'", alias)
	}

	dialer := proxy.New(info.Driver, info.Url, info.Username, info.Password, info.Secret)

	return dialer.Dial("tcp", string(target))
}

func (ProxyService) Create(info *model.CoreProxy) bool {
	if info.Url != "" {
		info.Url = lib.Encrypt(info.Url)
	}
	info.Username = lib.Encrypt(info.Username)
	if info.Password != "" {
		info.Password = lib.Encrypt(info.Password)
	}
	if info.Secret != "" {
		info.Secret = lib.Encrypt(info.Secret)
	}
	return model.DB().Create(info).RowsAffected > 0
}

func (ProxyService) Modify(info *model.CoreProxy) bool {
	value := *info
	if value.Url != "" {
		value.Url = lib.Encrypt(value.Url)
	}
	value.Username = lib.Encrypt(value.Username)
	if value.Password != "" {
		if !lib.IsHash(value.Password) {
			value.Password = lib.Encrypt(value.Password)
		} else {
			value.Password = ""
		}
	}
	if value.Secret != "" {
		if !lib.IsHash(value.Secret) {
			value.Secret = lib.Encrypt(value.Secret)
		} else {
			value.Secret = ""
		}
	}
	return 0 != info.ID && model.DB().Model(info).Where(&model.CoreProxy{ID: info.ID}).Update(&value).RowsAffected > 0
}

func (ProxyService) Status(id uint, status int) bool {
	cond := &model.CoreProxy{ID: id}
	return 0 != id && model.DB().Model(cond).Where(cond).Update("status", status).RowsAffected > 0
}

func (ProxyService) Page(start, end int) (count int, infos []model.CoreProxy) {
	model.DB().Model(&model.CoreProxy{}).Order("id desc").Count(&count).Offset(start).Limit(end).Find(&infos)

	for i := range infos {
		//if infos[i].Url != "" {
		//	infos[i].Url = lib.Decrypt(infos[i].Url)
		//}
		//if infos[i].Username != "" {
		//	infos[i].Username = lib.Decrypt(infos[i].Username)
		//}
		infos[i].Url = ""
		infos[i].Username = ""
		infos[i].Password = ""
		infos[i].Secret = ""
	}
	return
}

func (ProxyService) DetailById(id uint) *model.CoreProxy {
	if 0 == id {
		return nil
	}

	info := &model.CoreProxy{}
	if model.DB().Model(info).Where(&model.CoreProxy{ID: id}).Take(info).RowsAffected > 0 {
		if info.Url != "" {
			info.Url = lib.Decrypt(info.Url)
		}
		if info.Username != "" {
			info.Username = lib.Decrypt(info.Username)
		}
		if info.Password != "" {
			info.Password = lib.Hash(info.Password)
		}
		if info.Secret != "" {
			info.Secret = lib.Hash(info.Secret)
		}
		return info
	}

	return nil
}

func (ProxyService) InfoById(id uint) *model.CoreProxy {
	if 0 == id {
		return nil
	}

	info := &model.CoreProxy{}
	if model.DB().Model(info).Where(&model.CoreProxy{ID: id}).Take(info).RowsAffected > 0 {
		if info.Url != "" {
			info.Url = lib.Decrypt(info.Url)
		}
		if info.Username != "" {
			info.Username = lib.Decrypt(info.Username)
		}
		if info.Password != "" {
			info.Password = lib.Decrypt(info.Password)
		}
		if info.Secret != "" {
			info.Secret = lib.Decrypt(info.Secret)
		}
		return info
	}

	return nil
}

func (ProxyService) InfoByAlias(alias string) *model.CoreProxy {
	if "" == alias {
		return nil
	}

	info := &model.CoreProxy{}
	if model.DB().Model(info).Where(&model.CoreProxy{Alias: alias}).Take(info).RowsAffected > 0 {
		if info.Url != "" {
			info.Url = lib.Decrypt(info.Url)
		}
		if info.Username != "" {
			info.Username = lib.Decrypt(info.Username)
		}
		if info.Password != "" {
			info.Password = lib.Decrypt(info.Password)
		}
		if info.Secret != "" {
			info.Secret = lib.Decrypt(info.Secret)
		}
		return info
	}

	return nil
}
