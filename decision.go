package zen

// #include "zen_engine.h"
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

type decision struct {
	decisionPtr *C.ZenDecisionStruct
}

// newDecision: called internally by zen_engine only, cleanup should still be fired however.
func newDecision(decisionPtr *C.ZenDecisionStruct) Decision {
	return decision{
		decisionPtr: decisionPtr,
	}
}

func (decision decision) Evaluate(context any) (*EvaluationResponse, error) {
	return decision.EvaluateWithOpts(context, EvaluationOptions{})
}

func (decision decision) EvaluateWithOpts(context any, options EvaluationOptions) (*EvaluationResponse, error) {
	jsonData, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}

	cData := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cData))

	maxDepth := options.MaxDepth
	if maxDepth == 0 {
		maxDepth = 1
	}

	resultPtr := C.zen_decision_evaluate(decision.decisionPtr, cData, C.ZenEngineEvaluationOptions{
		trace:     C.bool(options.Trace),
		max_depth: C.uint8_t(maxDepth),
	})
	if resultPtr.error > 0 {
		var errorDetails string
		if resultPtr.details != nil {
			defer C.free(unsafe.Pointer(resultPtr.details))
			errorDetails = C.GoString(resultPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", resultPtr.error)
		}

		return nil, errors.New(errorDetails)
	}

	defer C.free(unsafe.Pointer(resultPtr.result))
	result := C.GoString(resultPtr.result)

	var response EvaluationResponse
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (decision decision) Dispose() {
	C.zen_decision_free(decision.decisionPtr)
}
