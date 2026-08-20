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
	initError              error
}

type EngineConfig struct {
	Loader            EngineLoader
	CustomNodeHandler CustomNodeHandler
}

type EngineLoader interface {
	engineLoader()
}

type LoaderConfig interface {
	EngineLoader
	cLoaderConfig() (C.ZenEngineLoaderConfig, func(), error)
}

type StaticLoader struct {
	Content map[string]json.RawMessage
}

func (StaticLoader) engineLoader() {}

func (config StaticLoader) cLoaderConfig() (C.ZenEngineLoaderConfig, func(), error) {
	data, err := json.Marshal(config.Content)
	if err != nil {
		return C.ZenEngineLoaderConfig{}, nil, err
	}

	content := C.CString(string(data))
	return C.ZenEngineLoaderConfig{
		kind:    C.ZenLoaderConfigKind_Static,
		content: content,
	}, func() { C.free(unsafe.Pointer(content)) }, nil
}

type FilesystemLoader struct {
	Path string
}

func (FilesystemLoader) engineLoader() {}

func (config FilesystemLoader) cLoaderConfig() (C.ZenEngineLoaderConfig, func(), error) {
	content := C.CString(config.Path)
	return C.ZenEngineLoaderConfig{
		kind:    C.ZenLoaderConfigKind_Filesystem,
		content: content,
	}, func() { C.free(unsafe.Pointer(content)) }, nil
}

type ZipLoader struct {
	Bytes []byte
}

func (ZipLoader) engineLoader() {}

func (config ZipLoader) cLoaderConfig() (C.ZenEngineLoaderConfig, func(), error) {
	if len(config.Bytes) == 0 {
		return C.ZenEngineLoaderConfig{}, nil, errors.New("zip loader config requires bytes")
	}

	bytes := C.CBytes(config.Bytes)
	return C.ZenEngineLoaderConfig{
		kind:      C.ZenLoaderConfigKind_Zip,
		bytes:     (*C.uint8_t)(bytes),
		bytes_len: C.uintptr_t(len(config.Bytes)),
	}, func() { C.free(bytes) }, nil
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

	var loaderCallback Loader
	var loaderConfig LoaderConfig

	switch loader := config.Loader.(type) {
	case nil:
	case Loader:
		loaderCallback = loader
	case LoaderConfig:
		loaderConfig = loader
	default:
		return engine{initError: fmt.Errorf("unsupported loader type %T", config.Loader)}
	}

	if config.CustomNodeHandler != nil {
		newEngine.customNodeHandler = cgo.NewHandle(wrapCustomNodeHandler(config.CustomNodeHandler))
		customNodeHandlerIdPtr = C.uintptr_t(newEngine.customNodeHandler)
		newEngine.customNodeHandlerIdPtr = &customNodeHandlerIdPtr
	}

	if loaderConfig != nil {
		cConfig, freeConfig, err := loaderConfig.cLoaderConfig()
		if err != nil {
			newEngine.disposeHandlers()
			return engine{initError: err}
		}
		defer freeConfig()

		resultPtr := C.zen_engine_new_golang_with_loader_config(cConfig, &customNodeHandlerIdPtr)
		if resultPtr.error > 0 {
			newEngine.disposeHandlers()

			var errorDetails string
			if resultPtr.details != nil {
				defer C.free(unsafe.Pointer(resultPtr.details))
				errorDetails = C.GoString(resultPtr.details)
			} else {
				errorDetails = fmt.Sprintf("Error code: %d", resultPtr.error)
			}

			return engine{initError: errors.New(errorDetails)}
		}

		newEngine.enginePtr = resultPtr.result
		return newEngine
	}

	if loaderCallback != nil {
		newEngine.loaderHandler = cgo.NewHandle(wrapLoader(loaderCallback))
		loaderHandlerIdPtr = C.uintptr_t(newEngine.loaderHandler)
		newEngine.loaderHandlerIdPtr = &loaderHandlerIdPtr
	}

	newEngine.enginePtr = C.zen_engine_new_golang(&loaderHandlerIdPtr, &customNodeHandlerIdPtr)
	return newEngine
}

func (engine engine) disposeHandlers() {
	if engine.loaderHandlerIdPtr != nil {
		engine.loaderHandler.Delete()
	}

	if engine.customNodeHandlerIdPtr != nil {
		engine.customNodeHandler.Delete()
	}
}

func (engine engine) Evaluate(key string, context any) (*EvaluationResponse, error) {
	return engine.EvaluateWithOpts(key, context, EvaluationOptions{})
}

func (engine engine) EvaluateWithOpts(key string, context any, options EvaluationOptions) (*EvaluationResponse, error) {
	if engine.initError != nil {
		return nil, engine.initError
	}

	jsonData, err := extractJsonFromAny(context)
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

func (engine engine) EvaluateBatch(requests []EvaluateBatchRequest) ([]EvaluateBatchResult, error) {
	return engine.EvaluateBatchWithOpts(requests, EvaluationOptions{})
}

func (engine engine) EvaluateBatchWithOpts(requests []EvaluateBatchRequest, options EvaluationOptions) ([]EvaluateBatchResult, error) {
	if engine.initError != nil {
		return nil, engine.initError
	}

	if len(requests) == 0 {
		return []EvaluateBatchResult{}, nil
	}

	cRequests := make([]C.ZenEngineEvaluateBatchRequest, len(requests))
	defer func() {
		for _, request := range cRequests {
			if request.key != nil {
				C.free(unsafe.Pointer(request.key))
			}
			if request.context != nil {
				C.free(unsafe.Pointer(request.context))
			}
		}
	}()

	for i, request := range requests {
		jsonData, err := extractJsonFromAny(request.Context)
		if err != nil {
			return nil, err
		}

		cRequests[i] = C.ZenEngineEvaluateBatchRequest{
			key:     C.CString(request.Key),
			context: C.CString(string(jsonData)),
		}
	}

	maxDepth := options.MaxDepth
	if maxDepth == 0 {
		maxDepth = 1
	}

	resultPtr := C.zen_engine_evaluate_batch(engine.enginePtr, &cRequests[0], C.uintptr_t(len(requests)), C.ZenEngineEvaluationOptions{
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

	var response []EvaluateBatchResult
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (engine engine) GetDecision(key string) (Decision, error) {
	if engine.initError != nil {
		return nil, engine.initError
	}

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
	if engine.initError != nil {
		return nil, engine.initError
	}

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
	if engine.enginePtr != nil {
		C.zen_engine_free(engine.enginePtr)
	}

	engine.disposeHandlers()
}
