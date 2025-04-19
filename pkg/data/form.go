package data

import (
	"fmt"
	"strings"
)

type FormTheme struct {
	FormClass         string // global form wrapper class
	LabelClass        string
	InputClass        string
	SelectClass       string
	CheckboxClass     string
	FieldWrapper      string // div class for each label+input
	RowWrapper        string // for arrays or grouped children
	FieldsetClass     string
	LegendClass       string
	AddButtonClass    string
	SubmitButtonClass string
	DeleteButtonClass string
	Repeat            int // default repeat count for arrays
}

func GenerateFormHTMLWithName(field Field, theme FormTheme, formName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<form method="POST" action="/submit" class="%s" name="%s">`, theme.FormClass, formName))
	renderFieldWithTheme(&sb, field, theme)
	sb.WriteString(fmt.Sprintf(`<button class="%s" type="submit">Submit</button></form>`, theme.SubmitButtonClass))
	return sb.String()
}

func GenerateFormHTML(field Field, theme FormTheme) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<form method="POST" action="/submit" class="%s">`, theme.FormClass))
	renderFieldWithTheme(&sb, field, theme)
	sb.WriteString(`<button type="submit">Submit</button></form>`)
	return sb.String()
}

func renderFieldWithTheme(sb *strings.Builder, f Field, theme FormTheme) {
	switch f.Type {
	case "object":
		if !f.IsRoot && f.Key != "" {
			sb.WriteString(fmt.Sprintf(`<fieldset class="%s">`, theme.FieldsetClass))
			sb.WriteString(fmt.Sprintf(`<legend class="%s">%s</legend>`, theme.LegendClass, f.Key))
		}
		for _, child := range f.Children {
			renderFieldWithTheme(sb, child, theme)
		}
		if !f.IsRoot && f.Key != "" {
			sb.WriteString(`</fieldset>`)
		}

	case "array":
		sb.WriteString(fmt.Sprintf(`<fieldset class="%s"><legend class="%s">%s</legend>`, theme.FieldsetClass, theme.LegendClass, f.Key))
		repeat := theme.Repeat
		if repeat == 0 {
			repeat = 5
		}
		for i := 0; i < repeat; i++ {
			sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.RowWrapper))
			for _, child := range f.Children {
				clone := updatePathForArrayIndex(child, i)
				renderFieldWithTheme(sb, clone, theme)
			}
			sb.WriteString(`</div><hr/>`)
		}
		sb.WriteString(`</fieldset>`)

	case "enum":
		sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		sb.WriteString(fmt.Sprintf(`<label class="%s">%s: <select name="%s" class="%s">`,
			theme.LabelClass, f.Key, f.Path, theme.SelectClass))
		for _, opt := range f.EnumVals {
			sb.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, opt, opt))
		}
		sb.WriteString(`</select></label><br/>`)
		sb.WriteString(`</div>`)

	case "bool":
		sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		sb.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="checkbox" name="%s" value="true" class="%s"/></label><br/>`,
			theme.LabelClass, f.Key, f.Path, theme.CheckboxClass))
		sb.WriteString(`</div>`)

	case "string":
		sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		sb.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="text" name="%s" class="%s"/></label><br/>`,
			theme.LabelClass, f.Key, f.Path, theme.InputClass))
		sb.WriteString(`</div>`)

	case "date", "number", "email":
		sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		sb.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="%s" name="%s" class="%s"/></label><br/>`,
			theme.LabelClass, f.Key, f.Type, f.Path, theme.InputClass))
		sb.WriteString(`</div>`)

	default:
		sb.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		sb.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="text" name="%s" class="%s"/></label><br/>`,
			theme.LabelClass, f.Key, f.Path, theme.InputClass))
		sb.WriteString(`</div>`)
	}
}

func updatePathForArrayIndex(f Field, index int) Field {
	newField := f

	// Remplacer UNIQUEMENT la première occurrence de ".0" ou suffix "0" dans le path
	// ex: matrix.0.0 → matrix.4.0  (mais pas matrix.4.4)
	updated := false
	parts := strings.Split(f.Path, ".")
	for i := 0; i < len(parts); i++ {
		if parts[i] == "0" && !updated {
			parts[i] = fmt.Sprintf("%d", index)
			updated = true
			break
		}
	}
	newField.Path = strings.Join(parts, ".")

	// Recurse on children
	newChildren := make([]Field, len(f.Children))
	for i, child := range f.Children {
		newChildren[i] = updatePathForArrayIndex(child, index)
	}
	newField.Children = newChildren

	return newField
}
