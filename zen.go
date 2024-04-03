package zen

import "encoding/json"

type EvaluationOptions struct {
	Trace    bool  `json:"trace"`
	MaxDepth uint8 `json:"maxDepth"`
}

type EvaluationResponse struct {
	Performance string           `json:"performance"`
	Result      json.RawMessage  `json:"result"`
	Trace       *json.RawMessage `json:"trace"`
}

type Engine interface {
	Evaluate(key string, context any) (*EvaluationResponse, error)
	EvaluateWithOpts(key string, context any, options EvaluationOptions) (*EvaluationResponse, error)
	GetDecision(key string) (Decision, error)
	CreateDecision(data []byte) (Decision, error)
	Dispose()
}

type Decision interface {
	Evaluate(context any) (*EvaluationResponse, error)
	EvaluateWithOpts(context any, options EvaluationOptions) (*EvaluationResponse, error)
	Dispose()
}
