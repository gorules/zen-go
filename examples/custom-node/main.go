package main

import (
	"embed"
	"fmt"
	"github.com/gorules/zen-go/v2"
	"github.com/gorules/zen-go/v2/examples/custom-node/nodes"
	"path"
)

//go:embed rules
var rulesFS embed.FS

func readTestFile(key string) ([]byte, error) {
	data, err := rulesFS.ReadFile(path.Join("rules", key))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: nodes.CustomNodeHandler})
	context := map[string]any{"a": 10}
	r, _ := engine.Evaluate("custom-node.json", context)

	fmt.Printf("[%s] Your result is: %s.\n", r.Performance, r.Result)
}
