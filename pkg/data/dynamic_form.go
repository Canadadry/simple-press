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
	b.WriteString(fmt.Sprintf(`<button class="%s" type="submit">Submit</button></form>`, theme.SubmitButtonClass))

	return b.String()
}

func renderDynamicField(b *strings.Builder, f Field, theme FormTheme) {
	switch f.Type {
	case "object":
		if f.IsRoot || f.Key == "" {
			for _, child := range f.Children {
				renderDynamicField(b, child, theme)
			}
			return
		}
		b.WriteString(fmt.Sprintf(`<fieldset class="%s"><legend class="%s">%s</legend>`,
			theme.FieldsetClass, theme.LegendClass, html.EscapeString(f.Key)))
		for _, child := range f.Children {
			renderDynamicField(b, child, theme)
		}
		b.WriteString(`</fieldset>`)

	case "array":
		b.WriteString(fmt.Sprintf(`<fieldset class="%s"><legend class="%s">%s</legend>`,
			theme.FieldsetClass, theme.LegendClass, html.EscapeString(f.Key)))

		containerID := "container-" + f.Path
		b.WriteString(fmt.Sprintf(`<div id="%s">`, containerID))

		for i := 0; i < 1; i++ { // une seule ligne initiale
			child := updatePathForArrayIndex(f.Children[0], i)
			b.WriteString(`<div>`)
			renderDynamicField(b, child, theme)
			b.WriteString(`</div>`)
		}
		b.WriteString(`</div>`)

		// Ajouter automatiquement le bouton "Add"
		b.WriteString(fmt.Sprintf(
			`<button type="button" class="%s" onclick="add_%s()">Add</button>`,
			theme.AddButtonClass, safeJSVar(f.Path),
		))

		// Template HTML pour duplication
		b.WriteString(fmt.Sprintf(`<template id="template-%s"><div>`, f.Path))
		templateField := replaceIndexInPath(f.Children[0], "__INDEX__")
		renderDynamicField(b, templateField, theme)
		b.WriteString(`</div></template>`)

		b.WriteString(`</fieldset>`)

	case "string", "number", "email", "date":
		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: `, theme.LabelClass, html.EscapeString(f.Key)))
		inputType := map[string]string{
			"string": "text", "number": "number", "email": "email", "date": "date",
		}[f.Type]
		b.WriteString(fmt.Sprintf(`<input type="%s" name="%s" class="%s"/>`,
			inputType, f.Path, theme.InputClass))
		b.WriteString(`</label><br/></div>`)

	case "bool":
		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: `, theme.LabelClass, html.EscapeString(f.Key)))
		b.WriteString(fmt.Sprintf(`<input type="checkbox" name="%s" value="true" class="%s"/>`,
			f.Path, theme.CheckboxClass))
		b.WriteString(`</label><br/></div>`)

	case "enum":
		b.WriteString(fmt.Sprintf(`<div class="%s">`, theme.FieldWrapper))
		b.WriteString(fmt.Sprintf(`<label class="%s">%s: `, theme.LabelClass, html.EscapeString(f.Key)))
		b.WriteString(fmt.Sprintf(`<select name="%s" class="%s">`, f.Path, theme.SelectClass))
		for _, opt := range f.EnumVals {
			b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`,
				html.EscapeString(opt), html.EscapeString(opt)))
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
