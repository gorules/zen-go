package nodes

import (
	"errors"
	"github.com/gorules/zen-go"
)

type NodeHandler interface {
	Handle(request zen.NodeRequest) (zen.NodeResponse, error)
}

var customNodes = map[string]NodeHandler{
	"add": addNode{},
	"mul": mulNode{},
	"sub": subNode{},
	"div": divNode{},
}

func CustomNodeHandler(request zen.NodeRequest) (zen.NodeResponse, error) {
	nodeHandler, ok := customNodes[request.Node.Kind]
	if !ok {
		return zen.NodeResponse{}, errors.New("component not found")
	}

	return nodeHandler.Handle(request)
}
