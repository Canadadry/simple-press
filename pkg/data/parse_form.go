package data

import (
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
			if _, ok := values[fullKey]; ok {
				result[key] = values.Get(fullKey)
			} else {
				result[key] = def[key]
			}

		case bool:
			if _, ok := values[fullKey]; ok {
				result[key] = values.Get(fullKey) == "true"
			} else {
				result[key] = def[key]
			}
		case int, int32, int64, float32, float64:
			if _, ok := values[fullKey]; ok {
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
			} else {
				result[key] = def[key]
			}
		}
	}

	return result, nil
}
