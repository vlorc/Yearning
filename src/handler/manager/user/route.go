package user

import (
	"Yearning-go/src/lib"
)

func SuperUserApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Put:    SuperFetchUser,
		Post:   ManageUserCreateOrEdit,
		Delete: SuperDeleteUser,
		Get:    ManageUserFetch,
	}
}
