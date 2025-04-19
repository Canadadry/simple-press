package data

import (
	"fmt"
	"html"
	"strings"
)

func GenerateFormDynamicHTMLWithName(root Field, theme FormTheme, formName string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`<form method="POST" action="/submit" name="%s" class="%s">`, html.EscapeString(formName), theme.FormClass))
	renderDynamicField(&b, root, theme)
	b.WriteString(`<button type="submit">Submit</button></form>`)

	return b.String()
}

func renderDynamicField(b *strings.Builder, f Field, theme FormTheme) {
	switch f.Type {
	case "object":
		if f.IsRoot {
			for _, child := range f.Children {
				renderDynamicField(b, child, theme)
			}
			return
		}

		if f.Key == "" {
			for _, child := range f.Children {
				renderDynamicField(b, child, theme)
			}
			return
		}

		b.WriteString(fmt.Sprintf(`<fieldset class="%s"><legend class="%s">%s</legend>`, theme.FieldsetClass, theme.LegendClass, html.EscapeString(f.Key)))
		for _, child := range f.Children {
			renderDynamicField(b, child, theme)
		}
		b.WriteString(`</fieldset>`)

	case "array":
		containerID := "container-" + f.Path
		templateID := "template-" + f.Path

		// Bloc initial visible : .0
		b.WriteString(fmt.Sprintf(`<fieldset class="%s"><legend class="%s">%s</legend>`, theme.FieldsetClass, theme.LegendClass, html.EscapeString(f.Key)))
		b.WriteString(fmt.Sprintf(`<div id="%s">`, containerID))

		child := updatePathForArrayIndex(f.Children[0], 0)
		b.WriteString(`<div>`)
		renderDynamicField(b, child, theme)
		b.WriteString(`</div>`)

		b.WriteString(`</div>`) // close container

		// Template avec __INDEX__
		templateChild := replaceIndexInPath(f.Children[0], "__INDEX__")
		b.WriteString(fmt.Sprintf(`<template id="%s">`, templateID))
		b.WriteString(`<div>`)
		renderDynamicField(b, templateChild, theme)
		b.WriteString(`</div>`)
		b.WriteString(`</template>`)

		b.WriteString(`</fieldset>`)

	case "string", "email", "date", "number":
		inputType := map[string]string{
			"string": "text",
			"email":  "email",
			"date":   "date",
			"number": "number",
		}[f.Type]

		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="%s" name="%s" class="%s"/></label><br/>`,
			theme.LabelClass,
			html.EscapeString(f.Key),
			inputType,
			html.EscapeString(f.Path),
			theme.InputClass,
		))
		b.WriteString(`</div>`)

	case "bool":
		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: <input type="checkbox" name="%s" value="true" class="%s"/></label><br/>`,
			theme.LabelClass,
			html.EscapeString(f.Key),
			html.EscapeString(f.Path),
			theme.CheckboxClass,
		))
		b.WriteString(`</div>`)

	case "enum":
		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: <select name="%s" class="%s">`,
			theme.LabelClass,
			html.EscapeString(f.Key),
			html.EscapeString(f.Path),
			theme.SelectClass,
		))
		for _, val := range f.EnumVals {
			b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, html.EscapeString(val), html.EscapeString(val)))
		}
		b.WriteString(`</select></label><br/></div>`)
	}
}

func replaceIndexInPath(f Field, indexToken string) Field {
	// remplace ".0" par .__INDEX__ ou similaire
	newPath := strings.Replace(f.Path, ".0", "."+indexToken, 1)
	if f.Path == "0" {
		newPath = indexToken
	}

	newF := f
	newF.Path = newPath

	newChildren := make([]Field, len(f.Children))
	for i, c := range f.Children {
		newChildren[i] = replaceIndexInPath(c, indexToken)
	}
	newF.Children = newChildren

	return newF
}

func GenerateDynamicJS(f Field) string {
	var b strings.Builder
	generateJSForField(&b, f)
	return b.String()
}

func generateJSForField(b *strings.Builder, f Field) {
	switch f.Type {
	case "array":
		safeID := safeJSVar(f.Path)

		// Déclare la variable d’index
		b.WriteString(fmt.Sprintf("let currentIndex_%s = 1;\n\n", safeID))

		// Fonction d'ajout
		b.WriteString(fmt.Sprintf("function add_%s() {\n", safeID))
		b.WriteString(fmt.Sprintf("  const template = document.getElementById(\"template-%s\").innerHTML;\n", f.Path))
		b.WriteString(fmt.Sprintf("  const container = document.getElementById(\"container-%s\");\n", f.Path))
		b.WriteString(fmt.Sprintf("  const html = template.replaceAll(\"__INDEX__\", currentIndex_%s);\n", safeID))
		b.WriteString("  const temp = document.createElement(\"div\");\n")
		b.WriteString("  temp.innerHTML = html;\n")
		b.WriteString("  container.appendChild(temp.firstElementChild);\n")
		b.WriteString(fmt.Sprintf("  currentIndex_%s++;\n", safeID))
		b.WriteString("}\n\n")

		// Gérer les enfants du champ du tableau (array of object ou of primitive)
		for _, c := range f.Children {
			generateJSForField(b, c)
		}

	case "object":
		for _, c := range f.Children {
			generateJSForField(b, c)
		}
	}
}

func safeJSVar(path string) string {
	return strings.ReplaceAll(path, ".", "_")
}
