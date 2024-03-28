package nodes

import (
	"errors"
	"github.com/gorules/zen-go"
)

type NodeHandler interface {
	Component() string
	Handle(request zen.NodeRequest) (zen.NodeResponse, error)
}

var customNodes = []NodeHandler{addNode{}}

func CustomNodeHandler(request zen.NodeRequest) (zen.NodeResponse, error) {
	for _, node := range customNodes {
		if node.Component() == request.Node.Component {
			return node.Handle(request)
		}
	}

	return zen.NodeResponse{}, errors.New("component not found")
}
