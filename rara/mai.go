package main

import (
	"github.com/gorules/zen-go"
	"github.com/gorules/zen-go/rara/nodes"
	"os"
	"path"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("test-data", key)
	return os.ReadFile(filePath)
}

func main() {
	engine := zen.NewEngine(zen.EngineConfig{Loader: readTestFile, CustomNodeHandler: nodes.CustomNodeHandler})
	decision, _ := engine.GetDecision("custom-node-aaa.json")

	for i := 0; i < 50; i++ {
		_, err := decision.EvaluateWithOpts(map[string]any{
			"a": 10,
			"b": 20,
			"c": 22,
		}, zen.EvaluationOptions{Trace: true})
		if err != nil {
			panic(err)
		}

		//r, _ := res.Trace.MarshalJSON()
		//ha := gjson.GetBytes(r, "138b3b11-ff46-450f-9704-3f3c712067b2.performance")
		//
		//rr, _ := res.Result.MarshalJSON()
		//
		//fmt.Printf("%+v; %+v; %+v\n", string(rr), ha.String(), res.Performance)
	}
}
