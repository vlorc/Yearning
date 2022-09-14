package settings

import (
	"Yearning-go/src/lib"
)

func SettingsApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Get:    SuperFetchSetting,
		Post:   SuperSaveSetting,
		Put:    SuperTestSetting,
		Delete: SuperDelOrder,
	}
}
