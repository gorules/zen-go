package zen

// #include "zen_engine.h"
import "C"

type Loader func(key string) ([]byte, error)

func wrapLoader(loader Loader) func(cKey *C.char) C.ZenDecisionLoaderResult {
	return func(cKey *C.char) C.ZenDecisionLoaderResult {
		key := C.GoString(cKey)
		content, err := loader(key)
		if err != nil {
			return C.ZenDecisionLoaderResult{
				content: nil,
				error:   C.CString(err.Error()),
			}
		}

		return C.ZenDecisionLoaderResult{
			content: C.CString(string(content)),
			error:   nil,
		}
	}
}
