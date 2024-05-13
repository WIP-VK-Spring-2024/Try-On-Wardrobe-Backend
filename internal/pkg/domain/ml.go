package domain

import "context"

//easyjson:json
type ModelHealthResponse struct {
	ConsumerCount   int
	MessageCount    int
	AvgResponseTime float32
	IsListening     bool
	ResponseTime    *float32
}

type MlModel interface {
	IsAvailable(model string, ctx context.Context) (bool, error)
}
