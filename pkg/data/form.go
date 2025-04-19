package data

import (
	"fmt"
	"strings"
)

// Main entrypoint
func GenerateFormHTML(root Field) string {
	var sb strings.Builder
	sb.WriteString(`<form method="POST" action="/submit">`)
	renderField(&sb, root)
	sb.WriteString(`<button type="submit">Submit</button></form>`)
	return sb.String()
}

func renderField(sb *strings.Builder, f Field) {
	switch f.Type {
	case "object":
		if !f.IsRoot && f.Key != "" {
			sb.WriteString(fmt.Sprintf("<fieldset><legend>%s</legend>", f.Key))
		}
		for _, child := range f.Children {
			renderField(sb, child)
		}
		if !f.IsRoot && f.Key != "" {
			sb.WriteString("</fieldset>")
		}

	case "array":
		sb.WriteString(fmt.Sprintf("<fieldset><legend>%s</legend>", f.Key))
		repeat := f.Repeat
		if repeat == 0 {
			repeat = 5
		}
		for i := 0; i < repeat; i++ {
			sb.WriteString("<div>")
			for _, child := range f.Children {
				clone := updatePathForArrayIndex(child, i)
				renderField(sb, clone)
			}
			sb.WriteString("</div><hr/>")
		}
		sb.WriteString("</fieldset>")

	case "enum":
		sb.WriteString(fmt.Sprintf(`<label>%s: <select name="%s">`, f.Key, f.Path))
		for _, opt := range f.EnumVals {
			sb.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, opt, opt))
		}
		sb.WriteString(`</select></label><br/>`)

	case "bool":
		sb.WriteString(fmt.Sprintf(`<label>%s: <input type="checkbox" name="%s" value="true" /></label><br/>`, f.Key, f.Path))

	case "string":
		sb.WriteString(fmt.Sprintf(`<label>%s: <input type="text" name="%s" /></label><br/>`, f.Key, f.Path))

	case "date", "number", "email":
		sb.WriteString(fmt.Sprintf(`<label>%s: <input type="%s" name="%s" /></label><br/>`, f.Key, f.Type, f.Path))

	default:
		sb.WriteString(fmt.Sprintf(`<label>%s: <input type="text" name="%s" /></label><br/>`, f.Key, f.Path))
	}
}

// Fixe les chemins des champs dans les tableaux : children.0.firstname → children.3.firstname
func updatePathForArrayIndex(f Field, index int) Field {
	parts := strings.Split(f.Path, ".")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "0" {
			parts[i] = fmt.Sprintf("%d", index)
			break
		}
	}
	f.Path = strings.Join(parts, ".")

	for i := range f.Children {
		f.Children[i] = updatePathForArrayIndex(f.Children[i], index)
	}
	return f
}
