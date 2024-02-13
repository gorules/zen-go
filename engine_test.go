package zen_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"sync"
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

func TestEngine_NewEngine(t *testing.T) {
	engineWithLoader := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engineWithLoader.Dispose()
	assert.NotNil(t, engineWithLoader)

	engineWithoutLoader := zen.NewEngine(zen.EngineConfig{})
	defer engineWithoutLoader.Dispose()
	assert.NotNil(t, engineWithoutLoader)
}

func TestEngine_Evaluate(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		var inputJson any
		err := json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := engine.Evaluate(data.file, inputJson)
		assert.NoError(t, err)
		assert.Nil(t, output.Trace)

		result, err := output.Result.MarshalJSON()
		assert.NoError(t, err)

		assert.Equal(t, data.outputJson, string(result))
	}
}

func TestEngine_EvaluateWithOpts(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		var inputJson any
		err := json.Unmarshal([]byte(data.inputJson), &inputJson)
		assert.NoError(t, err)

		output, err := engine.EvaluateWithOpts(data.file, inputJson, zen.EvaluationOptions{
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

func TestEngine_GetDecision(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	testData := prepareEvaluationTestData()
	for _, data := range testData {
		decision, err := engine.GetDecision(data.file)
		assert.NotNil(t, decision)
		assert.NoError(t, err)

		decision.Dispose()
	}
}

func TestEngine_CreateDecision(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	fileData, err := readTestFile("large.json")
	assert.NoError(t, err)

	decision, err := engine.CreateDecision(fileData)
	assert.NotNil(t, decision)
	assert.NoError(t, err)

	decision.Dispose()
}

func TestEngine_ErrorTransparency(t *testing.T) {
	errorStr := "Custom error"
	engine := zen.NewEngine(zen.EngineConfig{
		Loader: func(key string) ([]byte, error) {
			return nil, errors.New(errorStr)
		},
	})
	defer engine.Dispose()

	_, err := engine.Evaluate("myKey", nil)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "myKey")
	assert.ErrorContains(t, err, errorStr)
}

func TestEngine_EvaluateParallel(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile})
	defer engine.Dispose()

	type responseData struct {
		Output int `json:"output"`
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		current := i
		go func() {
			defer wg.Done()

			resp, err := engine.Evaluate("function.json", map[string]any{"input": current})
			assert.NoError(t, err)

			var respData responseData
			assert.NoError(t, json.Unmarshal(resp.Result, &respData))
			assert.Equal(t, current*2, respData.Output)
		}()
	}

	wg.Wait()
}
