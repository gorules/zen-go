package nodes

import (
	"github.com/gorules/zen-go"
)

type addNode struct {
}

func (a addNode) Component() string {
	return "add"
}

func (a addNode) Handle(request zen.NodeRequest) (zen.NodeResponse, error) {
	left, err := zen.GetNodeField[int](request, "left")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	right, err := zen.GetNodeField[int](request, "right")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	return zen.NodeResponse{
		Output: map[string]any{"sum": left + right},
	}, nil
}
