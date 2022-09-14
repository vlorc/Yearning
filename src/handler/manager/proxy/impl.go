package template

import "Yearning-go/src/model"

const (
	TEMPLATE_DELETE_SUCCESS = "权限组: %s 已删除"
	TEMPLATE_CREATE_SUCCESS = "%s代理已创建/编辑！"
	TEMPLATE_EDIT_SUCCESS   = "%s的代理已更新！"
	TEMPLATE_NOT_EXIST      = "大力不存在！"
)

type CommonProxyPost struct {
	Tp    string          `json:"tp"`
	Proxy model.CoreProxy `json:"proxy" json:"proxy"`
}
