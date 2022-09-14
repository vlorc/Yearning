package group

import (
	"Yearning-go/src/lib"
)

func GroupsApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:    SuperUserRuleMarge,
		Post:   SuperGroupUpdate,
		Put:    SuperGroup,
		Delete: SuperClearUserRule,
	}
}
