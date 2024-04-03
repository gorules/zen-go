package zen

// #include "zen_engine.h"
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

func EvaluateExpression[T any](expression string, context any) (T, error) {
	jsonData, err := extractJsonFromAny(context)
	if err != nil {
		var zero T
		return zero, err
	}

	expressionCString := C.CString(expression)
	defer C.free(unsafe.Pointer(expressionCString))

	contextCString := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(contextCString))

	resultPtr := C.zen_evaluate_expression(expressionCString, contextCString)
	if resultPtr.error > 0 {
		var errorDetails string
		if resultPtr.details != nil {
			defer C.free(unsafe.Pointer(resultPtr.details))
			errorDetails = C.GoString(resultPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", resultPtr.error)
		}

		var zero T
		return zero, errors.New(errorDetails)
	}

	defer C.free(unsafe.Pointer(resultPtr.result))
	resultJson := C.GoString(resultPtr.result)

	var result T
	if err := json.Unmarshal([]byte(resultJson), &result); err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}

func EvaluateUnaryExpression(expression string, context any) (bool, error) {
	jsonData, err := extractJsonFromAny(context)
	if err != nil {
		return false, err
	}

	expressionCString := C.CString(expression)
	defer C.free(unsafe.Pointer(expressionCString))

	contextCString := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(contextCString))

	resultPtr := C.zen_evaluate_unary_expression(expressionCString, contextCString)
	if resultPtr.error > 0 {
		var errorDetails string
		if resultPtr.details != nil {
			defer C.free(unsafe.Pointer(resultPtr.details))
			errorDetails = C.GoString(resultPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", resultPtr.error)
		}

		return false, errors.New(errorDetails)
	}

	isSuccess := int(*resultPtr.result)
	defer C.free(unsafe.Pointer(resultPtr.result))

	return isSuccess == 1, nil
}

func RenderTemplate[T any](template string, context any) (T, error) {
	jsonData, err := extractJsonFromAny(context)
	if err != nil {
		return *new(T), err
	}

	templateCString := C.CString(template)
	defer C.free(unsafe.Pointer(templateCString))

	contextCString := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(contextCString))

	resultPtr := C.zen_evaluate_template(templateCString, contextCString)
	if resultPtr.error > 0 {
		var errorDetails string
		if resultPtr.details != nil {
			defer C.free(unsafe.Pointer(resultPtr.details))
			errorDetails = C.GoString(resultPtr.details)
		} else {
			errorDetails = fmt.Sprintf("Error code: %d", resultPtr.error)
		}

		return *new(T), errors.New(errorDetails)
	}

	defer C.free(unsafe.Pointer(resultPtr.result))
	resultJson := C.GoString(resultPtr.result)

	var result T
	if err := json.Unmarshal([]byte(resultJson), &result); err != nil {
		return *new(T), err
	}

	return result, nil
}
