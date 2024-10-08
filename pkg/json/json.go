package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Parse(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func Write(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}
