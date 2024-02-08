package zen_test

import (
	"encoding/json"
	"github.com/gorules/zen-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecision_EvaluateWithOpts(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := engine.GetDecision(data.file)
		assert.NoError(t, err)

		var inputJson any
		err = json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := decision.Evaluate(inputJson)
		assert.NoError(t, err)
		assert.Nil(t, output.Trace)

		result, err := output.Result.MarshalJSON()
		assert.NoError(t, err)

		assert.Equal(t, data.outputJson, string(result))
		decision.Dispose()
	}
}

func TestDecision_Evaluate(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := engine.GetDecision(data.file)
		assert.NoError(t, err)

		var inputJson any
		err = json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := decision.EvaluateWithOpts(inputJson, zen.EvaluationOptions{
			Trace:    true,
			MaxDepth: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, output.Trace)

		result, err := output.Result.MarshalJSON()
		assert.NoError(t, err)

		assert.Equal(t, data.outputJson, string(result))
		decision.Dispose()
	}
}
