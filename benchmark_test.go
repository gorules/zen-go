package zen_test

import (
	"github.com/gorules/zen-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkEngine(b *testing.B) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile, CustomNodeHandler: customNodeHandler})
	defer engine.Dispose()

	context := map[string]any{"input": 5}

	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate("table.json", context)
	}
}

func BenchmarkDecision(b *testing.B) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile, CustomNodeHandler: customNodeHandler})
	defer engine.Dispose()

	decision, err := engine.GetDecision("table.json")
	require.NoError(b, err)
	defer decision.Dispose()

	context := map[string]any{"input": 5}

	for i := 0; i < b.N; i++ {
		_, _ = decision.Evaluate(context)
	}
}

func BenchmarkDecisionCustomNode(b *testing.B) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile, CustomNodeHandler: customNodeHandler})
	defer engine.Dispose()

	decision, err := engine.GetDecision("custom-node.json")
	require.NoError(b, err)
	defer decision.Dispose()

	context := map[string]any{"a": 5, "b": 10, "c": 15}

	for i := 0; i < b.N; i++ {
		_, _ = decision.Evaluate(context)
	}
}
