package audit

import (
	"Yearning-go/src/lib"
)

func AuditRestFulAPis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Post:   AuditOrderApis,
		Put: AuditOrRecordOrderFetchApis,
	}
}
