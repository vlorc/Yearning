package template

import (
	"Yearning-go/src/lib"
)

func TemplatesApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:    TemplateDetail,
		Post:   TemplateUpdate,
		Put:    TemplatePage,
	}
}
