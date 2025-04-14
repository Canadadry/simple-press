package view

import (
	"app/pkg/flash"
	"io"
)

var GoogleMapKey string

type ViewFunc func(io.Writer, func(string) string, flash.Message) error

func expandKind(in map[string]int) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(in))
	for name, value := range in {
		out = append(out, map[string]interface{}{
			"Value": value,
			"Name":  name,
		})
	}
	return out
}
