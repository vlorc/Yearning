package query

import (
	"Yearning-go/src/lib"
)

func AuditQueryRestFulAPis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Put:    AuditOrRecordQueryOrderFetchApis,
		Delete: QueryDeleteEmptyRecord,
		Post:   QueryHandlerSets,
	}
}
