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

type EvaluateBatchRequest struct {
	Key     string
	Context any
}

type EvaluateBatchResult struct {
	Success bool                `json:"success"`
	Data    *EvaluationResponse `json:"data,omitempty"`
	Error   json.RawMessage     `json:"error,omitempty"`
}

type Engine interface {
	Evaluate(key string, context any) (*EvaluationResponse, error)
	EvaluateWithOpts(key string, context any, options EvaluationOptions) (*EvaluationResponse, error)
	EvaluateBatch(requests []EvaluateBatchRequest) ([]EvaluateBatchResult, error)
	EvaluateBatchWithOpts(requests []EvaluateBatchRequest, options EvaluationOptions) ([]EvaluateBatchResult, error)
	GetDecision(key string) (Decision, error)
	CreateDecision(data []byte) (Decision, error)
	Dispose()
}

type Decision interface {
	Evaluate(context any) (*EvaluationResponse, error)
	EvaluateWithOpts(context any, options EvaluationOptions) (*EvaluationResponse, error)
	Dispose()
}
