package data

import (
	"fmt"
)

func ParseFormData(values, def map[string]any) (map[string]any, error) {
	result := map[string]any{}

	for key, subdef := range def {
		switch subdef_typed := subdef.(type) {
		case map[string]any:
			subvalues, ok := values[key].(map[string]any)
			if ok {
				var err error
				result[key], err = ParseFormData(subvalues, subdef_typed)
				if err != nil {
					return result, fmt.Errorf("%s: %w", key, err)
				}
			} else {
				result[key] = subdef_typed
			}

		case string:
			if v, ok := values[key].(string); ok {
				result[key] = v
			} else {
				result[key] = subdef_typed
			}

		case bool:
			if v, ok := values[key].(string); ok {
				result[key] = v == "true"
			} else {
				result[key] = subdef_typed
			}

		case int, int32, int64:
			if v, ok := values[key].(int); ok {
				result[key] = v
			} else if v64, ok := values[key].(int64); ok {
				result[key] = int(v64)
			} else if v, ok := values[key].(float64); ok {
				result[key] = int(v)
			} else {
				result[key] = subdef_typed
			}

		case float32, float64:
			if v, ok := values[key].(float64); ok {
				result[key] = v
			} else if v32, ok := values[key].(float32); ok {
				result[key] = float64(v32)
			} else if v, ok := values[key].(int); ok {
				result[key] = float64(v)
			} else {
				result[key] = subdef_typed
			}
		}
	}

	return result, nil
}
