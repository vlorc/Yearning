package db

import (
	"Yearning-go/src/lib"
)

func ManageDbApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Post:   ManageDBCreateOrEdit,
		Delete: SuperDeleteSource,
		Put:    SuperFetchSource,
	}
}
