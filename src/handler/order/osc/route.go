package osc

import (
	"Yearning-go/src/lib"
)

func AuditOSCFetchStateApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:    OscPercent,
		Delete: OscKill,
	}
}
