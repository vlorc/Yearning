package flow

import (
	"Yearning-go/src/lib"
)

func FlowRestApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:  GeneralAllSources,
		Post: FlowTplPostSourceTemplate,
		Put:  FlowTplEditSourceTemplateInfo,
	}
}
