package zen

import "encoding/json"

func extractJsonFromAny(data any) (json.RawMessage, error) {
	if d, ok := data.([]byte); ok {
		return d, nil
	}

	return json.Marshal(data)
}
