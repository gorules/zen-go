package zen_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"

	"github.com/gorules/zen-go"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("test-data", key)
	return os.ReadFile(filePath)
}

type evaluateTestData struct {
	file       string
	inputJson  string
	outputJson string
}

func prepareEvaluationTestData() map[string]evaluateTestData {
	return map[string]evaluateTestData{
		"table < 10": {
			file:       "table.json",
			inputJson:  `{"input":5}`,
			outputJson: `{"output":0}`,
		},
		"table > 10": {
			file:       "table.json",
			inputJson:  `{"input":15}`,
			outputJson: `{"output":10}`,
		},
		"function = 1": {
			file:       "function.json",
			inputJson:  `{"input":1}`,
			outputJson: `{"output":2}`,
		},
		"function = 5": {
			file:       "function.json",
			inputJson:  `{"input":5}`,
			outputJson: `{"output":10}`,
		},
		"function = 15": {
			file:       "function.json",
			inputJson:  `{"input":15}`,
			outputJson: `{"output":30}`,
		},
		"expression": {
			file:       "expression.json",
			inputJson:  `{"numbers": [1, 5, 15, 25],"firstName": "John","lastName": "Doe"}`,
			outputJson: `{"deep":{"nested":{"sum":46}},"fullName":"John Doe","largeNumbers":[15,25],"smallNumbers":[1,5]}`,
		},
	}
}

func TestZenEngine_NewEngine(t *testing.T) {
	zenEngineWithLoader := zen.NewEngine(readTestFile)
	defer zenEngineWithLoader.Dispose()
	assert.NotNil(t, zenEngineWithLoader)

	zenEngineWithoutLoader := zen.NewEngine(nil)
	defer zenEngineWithoutLoader.Dispose()
	assert.NotNil(t, zenEngineWithoutLoader)
}

func TestZenEngine_Evaluate(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		var inputJson any
		err := json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := zenEngine.Evaluate(data.file, inputJson)
		assert.NoError(t, err)
		assert.Nil(t, output.Trace)

		result, err := output.Result.MarshalJSON()
		assert.NoError(t, err)

		assert.Equal(t, data.outputJson, string(result))
	}
}

func TestZenEngine_EvaluateWithOpts(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		var inputJson any
		err := json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := zenEngine.EvaluateWithOpts(data.file, inputJson, zen.EvaluationOptions{
			Trace:    true,
			MaxDepth: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, output.Trace)

		result, err := output.Result.MarshalJSON()
		assert.NoError(t, err)

		assert.Equal(t, data.outputJson, string(result))
	}
}

func TestZenEngine_GetDecision(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := zenEngine.GetDecision(data.file)
		assert.NotNil(t, decision)
		assert.NoError(t, err)

		decision.Dispose()
	}
}

func TestZenEngine_CreateDecision(t *testing.T) {
	zenEngine := zen.NewEngine(readTestFile)
	defer zenEngine.Dispose()

	fileData, err := readTestFile("8k.json")
	assert.NoError(t, err)

	decision, err := zenEngine.CreateDecision(fileData)
	assert.NotNil(t, decision)
	assert.NoError(t, err)

	decision.Dispose()
}

func TestZenEngine_ErrorTransparency(t *testing.T) {
	errorStr := "Custom error from ZenEngine"
	zenEngine := zen.NewEngine(func(key string) ([]byte, error) {
		return nil, errors.New(errorStr)
	})

	_, err := zenEngine.Evaluate("", nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errorStr)
}
