package page

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/yuin/goldmark"
)

const BaseOf = "baseof.html"

type Data struct {
	Title         string
	Content       string
	Files         map[string]string
	Blocks        map[string]string
	ArticleBlocks map[string]map[string]any
}

func Render(w io.Writer, preview_data Data) error {
	if _, ok := preview_data.Files[BaseOf]; !ok {
		return fmt.Errorf("base template %s not defined", BaseOf)
	}
	funcMap := template.FuncMap{
		"markdownify": func(source string) template.HTML {
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(source), &buf); err != nil {
				return template.HTML(err.Error())
			}
			return template.HTML(buf.String())
		},
		"get_data": func(name string) map[string]any {
			return preview_data.ArticleBlocks[name]
		},
		"partial": func(name string, data map[string]any) (template.HTML, error) {
			content, ok := preview_data.Blocks[name]
			if !ok {
				return "", fmt.Errorf("unknown block %s", name)
			}
			buf := &bytes.Buffer{}
			tmpl, err := template.New("block").Parse(content)
			if err != nil {
				return "", fmt.Errorf("cannot parse block %s: %w", name, err)
			}
			err = tmpl.Execute(buf, data)
			if err != nil {
				return "", fmt.Errorf("cannot render block %s: %w", name, err)
			}
			return template.HTML(buf.String()), nil
		},
	}
	tmpl := template.New("").Funcs(funcMap)
	for name, content := range preview_data.Files {
		_, err := tmpl.New(name).Parse(content)
		if err != nil {
			return err
		}
	}
	type PageData struct {
		Title   string
		Blocks  map[string]map[string]any
		Content string
	}
	return tmpl.ExecuteTemplate(w, BaseOf, PageData{
		Content: preview_data.Content,
		Blocks:  preview_data.ArticleBlocks,
		Title:   preview_data.Title,
	})
}
