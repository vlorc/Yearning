package flow

type Step struct {
	Desc    string   `json:"desc"`
	Auditor []string `json:"auditor"`
	Type    int      `json:"type"`
}

type flowReq struct {
	Steps    []Step  `json:"steps"`
	Source   string `json:"source"`
	Relevant int    `json:"relevant"`
}

