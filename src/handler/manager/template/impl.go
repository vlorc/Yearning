package template

import "Yearning-go/src/model"

const (
	TEMPLATE_DELETE_SUCCESS = "权限组: %s 已删除"
	TEMPLATE_CREATE_SUCCESS = "%s模板已创建/编辑！"
	TEMPLATE_EDIT_SUCCESS   = "%s的模板已更新！"
	TEMPLATE_NOT_EXIST      = "模板不存在！"
)

type CommonTemplatePost struct {
	Tp       string             `json:"tp"`
	Template model.CoreTemplate `json:"template"`
}
