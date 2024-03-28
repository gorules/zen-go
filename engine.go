package zen

// #include "zen_engine.h"
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/cgo"
	"unsafe"
)

type engine struct {
	loaderHandler          cgo.Handle
	loaderHandlerIdPtr     *C.uintptr_t
	customNodeHandler      cgo.Handle
	customNodeHandlerIdPtr *C.uintptr_t
	enginePtr              *C.ZenEngineStruct
}

type EngineConfig struct {
	Loader            Loader
	CustomNodeHandler CustomNodeHandler
}

//export zen_engine_go_loader_callback
func zen_engine_go_loader_callback(h C.uintptr_t, key *C.char) C.ZenDecisionLoaderResult {
	fn := cgo.Handle(h).Value().(func(*C.char) C.ZenDecisionLoaderResult)
	return fn(key)
}

//export zen_engine_go_custom_node_callback
func zen_engine_go_custom_node_callback(h C.uintptr_t, request *C.char) C.ZenCustomNodeResult {
	fn := cgo.Handle(h).Value().(func(*C.char) C.ZenCustomNodeResult)
	return fn(request)
}

func NewEngine(config EngineConfig) Engine {
	var newEngine = engine{}
	var loaderHandlerIdPtr C.uintptr_t
	var customNodeHandlerIdPtr C.uintptr_t

	if config.Loader != nil {
		newEngine.loaderHandler = cgo.NewHandle(wrapLoader(config.Loader))
		loaderHandlerIdPtr = C.uintptr_t(newEngine.loaderHandler)
		newEngine.loaderHandlerIdPtr = &loaderHandlerIdPtr
	}

	if config.CustomNodeHandler != nil {
		newEngine.customNodeHandler = cgo.NewHandle(wrapCustomNodeHandler(config.CustomNodeHandler))
		customNodeHandlerIdPtr = C.uintptr_t(newEngine.customNodeHandler)
		newEngine.customNodeHandlerIdPtr = &customNodeHandlerIdPtr
	}

	newEngine.enginePtr = C.zen_engine_new_golang(&loaderHandlerIdPtr, &customNodeHandlerIdPtr)
	return newEngine
}

func (engine engine) Evaluate(key string, context any) (*EvaluationResponse, error) {
	return engine.EvaluateWithOpts(key, context, EvaluationOptions{})
}

func (engine engine) EvaluateWithOpts(key string, context any, options EvaluationOptions) (*EvaluationResponse, error) {
	jsonData, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cData := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cData))

	maxDepth := options.MaxDepth
	if maxDepth == 0 {
		maxDepth = 1
	}

	resultPtr := C.zen_engine_evaluate(engine.enginePtr, cKey, cData, C.ZenEngineEvaluationOptions{
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

func (engine engine) GetDecision(key string) (Decision, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	decisionPtr := C.zen_engine_get_decision(engine.enginePtr, cKey)
	if decisionPtr.error > 0 {
		var errorDetails string
		if decisionPtr.details != nil {
			defer C.free(unsafe.Pointer(decisionPtr.details))
			errorDetails = C.GoString(decisionPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", decisionPtr.error)
		}

		return nil, errors.New(errorDetails)
	}

	return newDecision(decisionPtr.result), nil
}

func (engine engine) CreateDecision(data []byte) (Decision, error) {
	cData := C.CString(string(data))
	defer C.free(unsafe.Pointer(cData))

	decisionPtr := C.zen_engine_create_decision(engine.enginePtr, cData)
	if decisionPtr.error > 0 {
		var errorDetails string
		if decisionPtr.details != nil {
			defer C.free(unsafe.Pointer(decisionPtr.details))
			errorDetails = C.GoString(decisionPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", decisionPtr.error)
		}

		return nil, errors.New(errorDetails)
	}

	return newDecision(decisionPtr.result), nil
}

func (engine engine) Dispose() {
	C.zen_engine_free(engine.enginePtr)

	if engine.loaderHandlerIdPtr != nil {
		engine.loaderHandler.Delete()
	}

	if engine.customNodeHandlerIdPtr != nil {
		engine.customNodeHandler.Delete()
	}
}
