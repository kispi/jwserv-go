package helper

import (
	"encoding/json"
)

// GetInputKeys GetInputKeys
func GetInputKeys(input []byte) []string {
	var objmap map[string]*json.RawMessage

	json.Unmarshal(input, &objmap)
	keys := make([]string, 0, len(objmap))
	for k := range objmap {
		keys = append(keys, k)
	}
	return keys
}
