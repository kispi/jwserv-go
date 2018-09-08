package helpers

import (
	"encoding/json"
	"strings"

	"../constants"
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

// MoveLimitToEnd moves 'limit' to the end so that subcount can be calculated before limit is applied.
func MoveLimitToEnd(queries []string) (result []string) {
	limit := ""
	for _, q := range queries {
		pair := strings.Split(q, "=")
		if pair[0] == constants.Limit {
			limit = q
		} else {
			result = append(result, q)
		}
	}
	if limit != "" {
		result = append(result, limit)
	}
	return
}
