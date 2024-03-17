package domain

type Result uint

const (
	ResultOk Result = iota + 1
	ResultRetry
	ResultDiscard
)
