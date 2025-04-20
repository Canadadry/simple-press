package data

import (
	"fmt"
	"sort"
	"strings"
)

type DynamicFormRenderer interface {
	BeginForm()
	EndForm()
	BeginFieldset(label string)
	EndFieldset()
	Input(label, name, inputType, value string)
	Checkbox(label, name string, checked bool)
}

func Render(data any, r DynamicFormRenderer) error {
	r.BeginForm()
	defer r.EndForm()
	return renderValue(data, "", "", r)
}

func renderValue(val any, key, path string, r DynamicFormRenderer) error {
	switch v := val.(type) {
	case map[string]any:
		if key != "" {
			r.BeginFieldset(key)
		}

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fullPath := joinPath(path, k)
			err := renderValue(v[k], k, fullPath, r)
			if err != nil {
				return err
			}
		}

		if key != "" {
			r.EndFieldset()
		}

	case string:
		r.Input(key, path, "text", v)
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		r.Input(key, path, "number", fmt.Sprintf("%v", v))
	case bool:
		r.Checkbox(key, path, v)
	default:
		return fmt.Errorf("at %v not handle value type %T", path, val)
	}
	return nil
}

func joinPath(parts ...string) string {
	nonEmpty := []string{}
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, ".")
}
