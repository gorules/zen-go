package nodes

import (
	"github.com/gorules/zen-go"
)

type addNode struct {
}

func (a addNode) Handle(request zen.NodeRequest) (zen.NodeResponse, error) {
	left, err := zen.GetNodeField[float64](request, "left")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	right, err := zen.GetNodeField[float64](request, "right")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	key, err := zen.GetNodeFieldRaw[string](request, "key")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	output := make(map[string]any)
	output[key] = left + right

	return zen.NodeResponse{Output: output}, nil
}
