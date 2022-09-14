package personal

import (
	"Yearning-go/src/lib"
)

func PersonalRestFulAPis()  lib.RestfulAPI{
	return lib.RestfulAPI{
		Post:    SQLReferToOrder,
		Put:    PersonalFetchOrderListOrProfile,
	}
}
