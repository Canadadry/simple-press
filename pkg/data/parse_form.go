package data

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

func ParseFormData(r *http.Request, def map[string]any) (map[string]any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return extractData(r.Form, def, "")
}

func extractData(values url.Values, def map[string]any, prefix string) (map[string]any, error) {
	result := map[string]any{}

	for key, val := range def {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch typed := val.(type) {
		case map[string]any:
			nested, err := extractData(values, typed, fullKey)
			if err != nil {
				return nil, err
			}
			result[key] = nested

		case string:
			result[key] = values.Get(fullKey)

		case bool:
			result[key] = values.Get(fullKey) == "true"

		case int, int32, int64, float32, float64:
			raw := values.Get(fullKey)
			if raw == "" {
				result[key] = 0.0
				continue
			}
			num, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return nil, err
			}
			result[key] = num

		default:
			// ignore arrays or unsupported types for now
			return nil, errors.New("unsupported type in schema: " + fullKey)
		}
	}

	return result, nil
}
