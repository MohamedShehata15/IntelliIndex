package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io"
)

// mustMarshalJSON marshals data to JSON and panics if there's an error
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}

// closeBody safely closes a response body and logs any errors
func closeBody(body io.ReadCloser) {
	if body == nil {
		return
	}

	err := body.Close()
	if err != nil {
		fmt.Printf("error closing response body: %v\n", err)
	}
}

// parseResponse is a generic helper function that parses an Elasticsearch response body into a given interface
func parseResponse(body io.ReadCloser, result interface{}) error {
	defer closeBody(body)

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	return nil
}
