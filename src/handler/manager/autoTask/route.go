package autoTask

import (
	"Yearning-go/src/lib"
)

func SuperAutoTaskApis() lib.RestfulAPI {
	return lib.RestfulAPI{
		Put:    SuperFetchAutoTaskList,
		Post:   SuperAutoTaskCreateOrEdit,
		Delete: SuperDeleteAutoTask,
	}
}
