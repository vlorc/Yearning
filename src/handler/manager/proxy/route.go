package template

import (
	"Yearning-go/src/lib"
)

func ProxyApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:  ProxyDetail,
		Post: ProxyUpdate,
		Put:  ProxyPage,
	}
}
