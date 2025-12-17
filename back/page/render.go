package page

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"sort"

	"github.com/yuin/goldmark"
)

const BaseOf = "baseof.html"

type Page struct {
	Slug        string
	Title       string
	Author      string
	Description string
}

type ArticleBlock struct {
	BlockName string
	Position  int
	Data      map[string]any
}

type Data struct {
	Title         string
	Content       string
	Files         map[string]string
	BlocksContent map[string]string
	ArticleBlocks []ArticleBlock
	PageFtecher   func(query string, offset, limit int) []Page
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
		"fetch": preview_data.PageFtecher,
		"partial": func(block ArticleBlock) (template.HTML, error) {
			content, ok := preview_data.BlocksContent[block.BlockName]
			if !ok {
				return "", fmt.Errorf("unknown block %s", block.BlockName)
			}
			buf := &bytes.Buffer{}
			tmpl, err := template.New("block").Parse(content)
			if err != nil {
				return "", fmt.Errorf("cannot parse block %s: %w", block.BlockName, err)
			}
			err = tmpl.Execute(buf, block.Data)
			if err != nil {
				return "", fmt.Errorf("cannot render block %s: %w", block.BlockName, err)
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
		Content string
		Blocks  []ArticleBlock
	}
	sort.Slice(preview_data.ArticleBlocks, func(i int, j int) bool {
		return preview_data.ArticleBlocks[i].Position < preview_data.ArticleBlocks[j].Position
	})
	return tmpl.ExecuteTemplate(w, BaseOf, PageData{
		Title:   preview_data.Title,
		Content: preview_data.Content,
		Blocks:  preview_data.ArticleBlocks,
	})
}
