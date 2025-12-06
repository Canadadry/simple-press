package httpquery

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHas(t *testing.T) {
	tests := []struct {
		url      string
		param    string
		expected bool
	}{
		{"http://example.org?toto=12", "toto", true},
		{"http://example.org?test=5", "test", true},
		{"http://example.org?test=5&test=6", "test", true},
		{"http://example.org?fake=5", "test", false},
	}
	for i, tt := range tests {
		req, err := http.NewRequest("GET", tt.url, nil)
		if err != nil {
			t.Fatalf("%d Cannot build query '%s' :%v", i, tt.url, err)
		}
		result := Has(req, tt.param)

		if result != tt.expected {
			t.Fatalf("%d Expect result %v but got %v", i, tt.expected, result)
		}

	}
}
func TestReadInt(t *testing.T) {
	tests := []struct {
		url          string
		path         string
		defaultValue int
		expected     int
	}{
		{"toto=12", "toto", 0, 12},
		{"test=5", "test", 0, 5},
		{"test=5&test=6", "test", 0, 0},
		{"test=a", "test", 0, 0},
		{"fake=5", "test", 0, 0},
	}
	for i, tt := range tests {
		base := "http://example.org?"
		req, err := http.NewRequest("GET", base+tt.url, nil)
		if err != nil {
			t.Fatalf("%d Cannot build query '%s%s' :%v", i, base, tt.url, err)
		}
		result := ReadInt(req, tt.path, tt.defaultValue)

		if result != tt.expected {
			t.Fatalf("%d Expect result %d but got %d", i, tt.expected, result)
		}

	}
}

func TestReadArray(t *testing.T) {
	tests := []struct {
		url          string
		path         string
		defaultValue map[string]string
		expected     map[string]string
	}{
		{"toto=12", "toto", nil, nil},
		{"toto[test]=12", "toto", nil, map[string]string{"test": "12"}},
		{"toto[test]=12&toto[test]=13", "toto", nil, map[string]string{"test": "12"}},
		{"toto[test]=12&tata[test]=13", "toto", nil, map[string]string{"test": "12"}},
		{"toto[test1]=12&toto[test2]=13", "toto", nil, map[string]string{"test1": "12", "test2": "13"}},
	}
	for i, tt := range tests {
		base := "http://example.org?"
		req, err := http.NewRequest("GET", base+tt.url, nil)
		if err != nil {
			t.Fatalf("%d Cannot build query '%s%s' :%v", i, base, tt.url, err)
		}
		result := ReadArray(req, tt.path, tt.defaultValue)

		if !reflect.DeepEqual(result, tt.expected) {
			t.Fatalf("%d Expect result %#v but got %#v", i, tt.expected, result)
		}

	}
}
