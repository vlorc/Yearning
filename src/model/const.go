package model

const (
	// 审核中
	ORDER_STATUS_AUDIT = 2
	// 驳回
	ORDER_STATUS_REJECT = 0
	// 已执行
	ORDER_STATUS_EXEC = 1
	// 执行失败
	ORDER_STATUS_EXEC_FAILED = 4
	// 待执行
	ORDER_STATUS_EXEC_READY = 5
)
