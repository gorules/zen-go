package zen_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/gorules/zen-go/v2"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("test-data", key)
	return os.ReadFile(filePath)
}

func customNodeHandler(request zen.NodeRequest) (zen.NodeResponse, error) {
	if request.Node.Kind != "sum" {
		return zen.NodeResponse{}, errors.New("unknown component")
	}

	a, err := zen.GetNodeField[int](request, "a")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	b, err := zen.GetNodeField[int](request, "b")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	key, err := zen.GetNodeFieldRaw[string](request, "key")
	if err != nil {
		return zen.NodeResponse{}, err
	}

	output := make(map[string]any)
	output[key] = a + b

	return zen.NodeResponse{Output: output}, nil
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
		"customNode": {
			file:       "custom-node.json",
			inputJson:  `{"a": 5, "b": 10, "c": 15}`,
			outputJson: `{"sum":30}`,
		},
	}
}

func TestEngine_NewEngine(t *testing.T) {
	engineWithLoader := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
	defer engineWithLoader.Dispose()
	assert.NotNil(t, engineWithLoader)

	engineWithoutLoader := zen.NewEngine(zen.EngineConfig{})
	defer engineWithoutLoader.Dispose()
	assert.NotNil(t, engineWithoutLoader)
}

func TestEngine_Evaluate(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
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

		assert.JSONEq(t, data.outputJson, string(result))
	}
}

func TestEngine_EvaluateWithOpts(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
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

		assert.JSONEq(t, data.outputJson, string(result))
	}
}

func TestEngine_GetDecision(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
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
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
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
		Loader: zen.Loader(func(key string) ([]byte, error) {
			return nil, errors.New(errorStr)
		}),
	})
	defer engine.Dispose()

	_, err := engine.Evaluate("myKey", nil)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "myKey")
	assert.ErrorContains(t, err, errorStr)
}

func TestEngine_EvaluateParallel(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile), CustomNodeHandler: customNodeHandler})
	defer engine.Dispose()

	type responseData struct {
		Output int `json:"output"`
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
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

func TestEngine_NewEngineWithStaticLoader(t *testing.T) {
	tableContent, err := readTestFile("table.json")
	assert.NoError(t, err)

	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.StaticLoader{
		Content: map[string]json.RawMessage{"table.json": tableContent},
	}})
	defer engine.Dispose()

	output, err := engine.Evaluate("table.json", map[string]any{"input": 15})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"output":10}`, string(output.Result))

	_, err = engine.Evaluate("missing.json", map[string]any{})
	assert.Error(t, err)
}

func TestEngine_NewEngineWithFilesystemLoader(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.FilesystemLoader{Path: "test-data"}})
	defer engine.Dispose()

	output, err := engine.Evaluate("table.json", map[string]any{"input": 15})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"output":10}`, string(output.Result))
}

func TestEngine_NewEngineWithZipLoader(t *testing.T) {
	tableContent, err := readTestFile("table.json")
	assert.NoError(t, err)

	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, err := writer.Create("table.json")
	assert.NoError(t, err)
	_, err = file.Write(tableContent)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.ZipLoader{Bytes: buffer.Bytes()}})
	defer engine.Dispose()

	output, err := engine.Evaluate("table.json", map[string]any{"input": 15})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"output":10}`, string(output.Result))
}

func TestEngine_NewEngineWithLoaderConfigCustomNode(t *testing.T) {
	customNodeContent, err := readTestFile("custom-node.json")
	assert.NoError(t, err)

	engine := zen.NewEngine(zen.EngineConfig{
		Loader: zen.StaticLoader{
			Content: map[string]json.RawMessage{"custom-node.json": customNodeContent},
		},
		CustomNodeHandler: customNodeHandler,
	})
	defer engine.Dispose()

	output, err := engine.Evaluate("custom-node.json", map[string]any{"a": 5, "b": 10, "c": 15})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sum":30}`, string(output.Result))
}

func TestEngine_NewEngineWithInvalidZipLoader(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.ZipLoader{Bytes: []byte{1, 2, 3, 4}}})
	defer engine.Dispose()

	_, err := engine.Evaluate("table.json", map[string]any{"input": 15})
	assert.Error(t, err)
}

func TestEngine_NewEngineWithTypedLoaderCallback(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile)})
	defer engine.Dispose()

	output, err := engine.Evaluate("table.json", map[string]any{"input": 15})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"output":10}`, string(output.Result))
}

func TestEngine_EvaluateBatch(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.FilesystemLoader{Path: "test-data"}})
	defer engine.Dispose()

	results, err := engine.EvaluateBatch([]zen.EvaluateBatchRequest{
		{Key: "table.json", Context: map[string]any{"input": 15}},
		{Key: "missing.json", Context: map[string]any{}},
		{Key: "table.json", Context: map[string]any{"input": 5}},
	})
	assert.NoError(t, err)
	assert.Len(t, results, 3)

	assert.True(t, results[0].Success)
	assert.JSONEq(t, `{"output":10}`, string(results[0].Data.Result))

	assert.False(t, results[1].Success)
	assert.NotEmpty(t, results[1].Error)

	assert.True(t, results[2].Success)
	assert.JSONEq(t, `{"output":0}`, string(results[2].Data.Result))
}

func TestEngine_EvaluateBatchEmpty(t *testing.T) {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.FilesystemLoader{Path: "test-data"}})
	defer engine.Dispose()

	results, err := engine.EvaluateBatch(nil)
	assert.NoError(t, err)
	assert.Empty(t, results)
}
