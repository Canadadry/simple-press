package maputil

import (
	"reflect"
	"testing"
)

func TestFlattern(t *testing.T) {
	tests := map[string]struct {
		input     map[string]interface{}
		separator string
		expected  map[string]interface{}
	}{
		"nil input": {
			nil,
			".",
			map[string]interface{}{},
		},
		"single nil value": {
			map[string]interface{}{"test1": nil},
			".",
			map[string]interface{}{"test1": nil},
		},
		"single flat key-value pair": {
			map[string]interface{}{"test1": "test2"},
			".",
			map[string]interface{}{"test1": "test2"},
		},
		"nested one level": {
			map[string]interface{}{"test1": map[string]interface{}{"test2": "test3"}},
			".",
			map[string]interface{}{"test1.test2": "test3"},
		},
		"nested multiple levels": {
			map[string]interface{}{
				"test1": map[string]interface{}{
					"test2": map[string]interface{}{
						"test3": "test4",
					},
				}},
			".",
			map[string]interface{}{"test1.test2.test3": "test4"},
		},
		"array of nested maps": {
			map[string]interface{}{
				"test1": []interface{}{
					map[string]interface{}{
						"test2": map[string]interface{}{
							"test3": "test4",
						},
					},
					map[string]interface{}{
						"test5": map[string]interface{}{
							"test6": "test7",
						},
					},
				},
			},
			".",
			map[string]interface{}{
				"test1.0.test2.test3": "test4",
				"test1.1.test5.test6": "test7",
			},
		},
		"array of nested int": {
			map[string]interface{}{
				"test1": []int{0, 1, 2},
			},
			".",
			map[string]interface{}{
				"test1.0": 0,
				"test1.1": 1,
				"test1.2": 2,
			},
		},
		"another array of nested map": {
			map[string]interface{}{
				"array_field": []map[string]interface{}{
					{"sub_field": "field_0"},
					{"sub_field": "field_1"},
					{"sub_field": "field_2"},
				},
			}, ".",
			map[string]interface{}{
				"array_field.0.sub_field": "field_0",
				"array_field.1.sub_field": "field_1",
				"array_field.2.sub_field": "field_2",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Flattern(tt.input, tt.separator)
			if !reflect.DeepEqual(tt.expected, result) {
				t.Fatalf("got %v expected %v", result, tt.expected)
			}
		})
	}
}

func TestExpand(t *testing.T) {
	tests := map[string]struct {
		input     interface{}
		separator string
		expected  map[string]interface{}
	}{
		"nil input": {
			nil,
			".",
			nil,
		},
		"single nil value": {
			map[string]interface{}{"test1": nil},
			".",
			map[string]interface{}{"test1": nil},
		},
		"flat key-value": {
			map[string]interface{}{"test1": "test2"},
			".",
			map[string]interface{}{"test1": "test2"},
		},
		"one level nested map": {
			map[string]interface{}{"test1.test2": "test3"},
			".",
			map[string]interface{}{"test1": map[string]interface{}{"test2": "test3"}},
		},
		"multiple levels nested map": {
			map[string]interface{}{"test1.test2.test3": "test4"},
			".",
			map[string]interface{}{
				"test1": map[string]interface{}{
					"test2": map[string]interface{}{
						"test3": "test4",
					},
				},
			},
		},
		"array of ints": {
			map[string]interface{}{
				"numbers.0": 1,
				"numbers.1": 2,
				"numbers.2": 3,
			},
			".",
			map[string]interface{}{
				"numbers": []interface{}{1, 2, 3},
			},
		},
		"validator error format": {
			map[string][]string{
				"numbers.0.name": []string{"test.0"},
				"numbers.1.name": []string{"test.1"},
				"numbers.2.name": []string{"test.2"},
			},
			".",
			map[string]interface{}{
				"numbers": []interface{}{
					map[string]interface{}{"name": []string{"test.0"}},
					map[string]interface{}{"name": []string{"test.1"}},
					map[string]interface{}{"name": []string{"test.2"}},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Expand(tt.input, tt.separator)
			if !reflect.DeepEqual(tt.expected, result) {
				t.Fatalf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func testGetSortedKeys(t *testing.T) {
	tests := map[string]struct {
		input    map[string]interface{}
		expected []string
	}{
		"all empty types": {
			map[string]interface{}{},
			[]string{},
		},
		"not empty values": {
			map[string]interface{}{
				"integer": 1,
				"float":   1.0,
				"string":  "test",
				"slice":   []string{"test"},
				"map": map[string]interface{}{
					"integer": 1,
					"float":   1.0,
					"string":  "test",
				},
				"bool": true,
				"struct": struct {
					test string
				}{
					test: "hello",
				},
			},
			[]string{
				"bool",
				"float",
				"integer",
				"map",
				"slice",
				"string",
				"struct",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			keys := GetSortedKeys[interface{}](tt.input)
			if !reflect.DeepEqual(keys, tt.expected) {
				t.Fatalf("failed : expected data %#v, got %#v", tt.expected, tt.input)
			}
		})
	}
}
