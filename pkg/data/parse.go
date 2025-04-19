package data

import (
	"sort"
	"strings"
)

type Field struct {
	Key      string
	Path     string
	Type     string
	EnumVals []string
	Children []Field
	Repeat   int
	IsRoot   bool
}

func Parse(input any, isRoot bool) Field {
	return parseInternal(input, isRoot, "")
}

func parseInternal(input any, isRoot bool, currentPath string) Field {
	switch v := input.(type) {

	case map[string]any:
		children := []Field{}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			child := parseField(k, v[k], currentPath)
			children = append(children, child)
		}
		return Field{
			IsRoot:   isRoot,
			Type:     "object",
			Path:     currentPath,
			Children: children,
		}

	case []any:
		f := Field{
			IsRoot: isRoot,
			Type:   "array",
			Path:   currentPath,
			Repeat: 5,
		}
		if len(v) > 0 {
			child := v[0]
			childPath := currentPath
			if currentPath != "" {
				childPath = currentPath + ".0"
			} else {
				childPath = "0"
			}
			childField := parseField("", child, childPath)

			// si objet, on le force comme tel
			if m, ok := child.(map[string]any); ok {
				childField = parseInternal(m, false, childPath)
				childField.Type = "object"
			}

			f.Children = []Field{childField}
		}
		return f
	}

	return Field{
		IsRoot: isRoot,
		Type:   "string",
		Path:   currentPath,
	}
}

func parseField(key string, val any, parentPath string) Field {
	var path string
	if parentPath == "" {
		path = key
	} else if key == "" {
		path = parentPath
	} else {
		path = parentPath + "." + key
	}

	switch t := val.(type) {
	case string:
		if strings.HasPrefix(t, "enum:") {
			enumVals := strings.Split(strings.TrimPrefix(t, "enum:"), ";")
			return Field{Key: key, Path: path, Type: "enum", EnumVals: enumVals}
		}
		return Field{Key: key, Path: path, Type: t}

	case map[string]any:
		obj := parseInternal(t, false, path)
		obj.Key = key
		return obj

	case []any:
		arr := Field{
			Key:    key,
			Path:   path,
			Type:   "array",
			Repeat: 5,
		}
		if len(t) > 0 {
			child := parseField("", t[0], path+".0")

			// si objet, le forcer comme type object avec path ajusté
			if m, ok := t[0].(map[string]any); ok {
				child = parseInternal(m, false, path+".0")
				child.Type = "object"
			}

			arr.Children = []Field{child}
		}
		return arr
	}

	return Field{Key: key, Path: path, Type: "string"}
}
