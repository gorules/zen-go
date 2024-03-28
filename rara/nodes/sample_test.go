package nodes

import (
	"github.com/gorules/zen-go"
	"os"
	"path"
	"testing"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("../..", "test-data", key)
	return os.ReadFile(filePath)
}

func BenchmarkCustomNodeHandler(b *testing.B) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile, CustomNodeHandler: CustomNodeHandler})
	decision, _ := engine.GetDecision("custom-node-aaa.json")
	context := map[string]any{
		"a": 10,
		"b": 20,
		"c": 22,
	}

	for n := 0; n < b.N; n++ {
		_, _ = decision.Evaluate(context)
	}
}
