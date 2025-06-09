package elasticsearch

import (
	"encoding/json"
	"fmt"
)

// mustMarshalJSON marshals data to JSON and panics if there's an error
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}
