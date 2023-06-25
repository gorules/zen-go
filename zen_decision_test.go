package zen_test

import (
	"encoding/json"
	"github.com/gorules/zen-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZenDecision_EvaluateWithOpts(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := zenEngine.GetDecision(data.file)
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

func TestZenDecision_Evaluate(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := zenEngine.GetDecision(data.file)
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
