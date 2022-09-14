package roles

import "Yearning-go/src/lib"

func RolesApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Put:  SuperFetchRoles,
		Post: SuperSaveRoles,
	}
}
