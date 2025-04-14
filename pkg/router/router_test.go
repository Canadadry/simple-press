package router

import (
	"testing"
)

func TestGetPatternFromURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		replace map[string]string
		result  string
	}{
		{
			name:    "no change with slash at start and end",
			url:     "/part1/part2/part3/part4/",
			replace: map[string]string{},
			result:  "/part1/part2/part3/part4/",
		},
		{
			name:    "no change without slash at start and end",
			url:     "part1/part2/part3/part4",
			replace: map[string]string{},
			result:  "part1/part2/part3/part4",
		},
		{
			name:    "digit replacement",
			url:     "/part1/part2/part3/1234567890",
			replace: map[string]string{},
			result:  "/part1/part2/part3/:digit",
		},
		{
			name:    "uuid replacement",
			url:     "/part1/part2/part3/23c72f7b-0bfb-47cc-8a5d-172ccb873f78",
			replace: map[string]string{},
			result:  "/part1/part2/part3/:uuid",
		},
		{
			name: "custom regexp replacement 1",
			url:  "/part1/part2/part3/part4",
			replace: map[string]string{
				"file/:filename": "part3/part4$",
			},
			result: "/part1/part2/file/:filename",
		},
		{
			name: "custom regexp replacement 2",
			url:  "/part1/part2/part3/part4",
			replace: map[string]string{
				"${1}file/:filename${2}": "^(/part1/)part2/part3(/part4)$",
			},
			result: "/part1/file/:filename/part4",
		},
		{
			name: "custom regexp and digit replacement",
			url:  "/part1/part2/part3/1234567890",
			replace: map[string]string{
				"/${1}/file/:filename/${2}": "^/([a-z0-9]+)/part2/part3/([0-9]+)$",
			},
			result: "/part1/file/:filename/:digit",
		},
		{
			name: "custom regexp and uuid replacement",
			url:  "/part1/part2/part3/23c72f7b-0bfb-47cc-8a5d-172ccb873f78",
			replace: map[string]string{
				"file/:filename/${1}": "part2/[a-z0-9]+/([a-z0-9\\-]+)$",
			},
			result: "/part1/file/:filename/:uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetPatternFromURL(tt.url, tt.replace)
			if tt.result != r {
				t.Fatalf("expected %s got  %s", tt.result, r)
			}
		})
	}

}
