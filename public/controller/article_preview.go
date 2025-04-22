package controller

import (
	"app/pkg/router"
	"app/public/repository"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/yuin/goldmark"
)

func (c *Controller) GetArticlePreview(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	if slug == "" {
		slug = "index"
	}
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return c.GetFile(w, r)
	}
	baseTemplates, err := c.Repository.SelectAllTemplate(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base template : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseTemplates {
		files[l.Name] = l.Content
	}
	pageLayout, ok, err := c.Repository.SelectLayoutByID(r.Context(), a.LayoutID)
	if err != nil {
		return fmt.Errorf("cannot select page layout %d : %w", a.LayoutID, err)
	}
	if !ok {
		return fmt.Errorf("cannot found page layout %d : %w", a.LayoutID, err)
	}
	files[pageLayout.Name] = pageLayout.Content

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := map[string]string{}
	for _, b := range blocks {
		blockSelector[b.Name] = b.Content
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), a.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := map[string]map[string]any{}
	for _, p := range blockDatas {
		blockDataView[p.BlockName] = p.Data
	}
	return renderPreview(w, files, blockSelector, blockDataView, a)
}

func renderPreview(w io.Writer, files map[string]string, blocks map[string]string, ArticleBlocks map[string]map[string]any, pageData repository.Article) error {
	const baseTemplate = "baseof.html"
	if _, ok := files[baseTemplate]; !ok {
		return fmt.Errorf("base template %s not defined", baseTemplate)
	}
	funcMap := template.FuncMap{
		"markdownify": func(source string) template.HTML {
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(source), &buf); err != nil {
				return template.HTML(err.Error())
			}
			return template.HTML(buf.String())
		},
		"partial": func(name string, data map[string]any) (template.HTML, error) {
			content, ok := blocks[name]
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
	for name, content := range files {
		if name == baseTemplate {
			continue
		}
		_, err := tmpl.New(name).Parse(content)
		if err != nil {
			return err
		}
	}
	_, err := tmpl.New(baseTemplate).Parse(files[baseTemplate])
	if err != nil {
		return err
	}
	type PageData struct {
		Title   string
		Blocks  map[string]map[string]any
		Content string
	}
	return tmpl.ExecuteTemplate(w, baseTemplate, PageData{
		Content: pageData.Content,
		Blocks:  ArticleBlocks,
		Title:   pageData.Title,
	})
}
