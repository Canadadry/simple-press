package maputil

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func Flattern(data any, separator string) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive("", data, separator, result)
	return result
}

func flattenRecursive(prefix string, value any, separator string, result map[string]interface{}) {
	if value == nil {
		if prefix != "" {
			result[prefix] = nil
		}
		return
	}

	val := reflect.ValueOf(value)
	kind := val.Kind()

	switch kind {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			strKey := fmt.Sprintf("%v", key)
			fullKey := strKey
			if prefix != "" {
				fullKey = prefix + separator + strKey
			}
			flattenRecursive(fullKey, val.MapIndex(key).Interface(), separator, result)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			fullKey := fmt.Sprintf("%s%s%d", prefix, separator, i)
			flattenRecursive(fullKey, val.Index(i).Interface(), separator, result)
		}
	default:
		if prefix != "" {
			result[prefix] = value
		}
	}
}

func isMap(v any) bool {
	typeOfValue := reflect.TypeOf(v)
	return typeOfValue.Kind() == reflect.Map
}

// Expand reconstruit une structure imbriquée (maps/slices) à partir d'une map aplatie avec des clés "a.b.0.c"
func Expand(flat interface{}, separator string) map[string]interface{} {
	if flat == nil {
		return nil
	}

	// Conversion dynamique vers map[string]interface{}
	flatMap := make(map[string]interface{})
	val := reflect.ValueOf(flat)
	if val.Kind() != reflect.Map {
		panic("Expand: input must be a map")
	}

	for _, key := range val.MapKeys() {
		strKey := fmt.Sprintf("%v", key.Interface())
		flatMap[strKey] = val.MapIndex(key).Interface()
	}

	root := make(map[string]interface{})

	for flatKey, value := range flatMap {
		keys := strings.Split(flatKey, separator)
		current := root

		for i := 0; i < len(keys); i++ {
			key := keys[i]
			isLast := i == len(keys)-1

			if isLast {
				current[key] = value
				break
			}

			if _, ok := current[key]; !ok {
				current[key] = map[string]interface{}{}
			}

			next, ok := current[key].(map[string]interface{})
			if !ok {
				newMap := map[string]interface{}{}
				current[key] = newMap
				next = newMap
			}

			current = next
		}
	}

	normalized := normalizeSlices(root)
	result, ok := normalized.(map[string]interface{})
	if !ok {
		panic("Expand: expected root to be map[string]interface{} after normalization")
	}
	return result
}

func normalizeSlices(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		// Si toutes les clés sont des indices numériques, on convertit en slice
		if isArrayLike(val) {
			maxIdx := -1
			for k := range val {
				idx, _ := strconv.Atoi(k)
				if idx > maxIdx {
					maxIdx = idx
				}
			}
			slice := make([]interface{}, maxIdx+1)
			for k, sub := range val {
				idx, _ := strconv.Atoi(k)
				slice[idx] = normalizeSlices(sub)
			}
			return slice
		}
		for k, sub := range val {
			val[k] = normalizeSlices(sub)
		}
		return val
	case []interface{}:
		for i, elem := range val {
			val[i] = normalizeSlices(elem)
		}
	}
	return v
}

func isArrayLike(m map[string]interface{}) bool {
	if len(m) == 0 {
		return false
	}
	for k := range m {
		if _, err := strconv.Atoi(k); err != nil {
			return false
		}
	}
	return true
}

func GetSortedKeys[T any](data map[string]T) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
