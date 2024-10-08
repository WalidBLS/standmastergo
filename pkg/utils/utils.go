package utils

import (
	"fmt"
	"net/http"
)

func GetIntFromMap(input map[string]interface{}, key string) (int, error) {
	value, ok := input[key]
	if !ok || value == nil {
		return 0, fmt.Errorf("%s is missing or nil", key)
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("%s is not a valid number", key)
	}

	return int(floatValue), nil
}

func GetQueryParams(r *http.Request) map[string]interface{} {
	query := r.URL.Query()
	params := map[string]interface{}{}

	for key, value := range query {
		if len(value) == 1 {
			params[key] = value[0]
		} else {
			params[key] = value
		}
	}

	return params
}
