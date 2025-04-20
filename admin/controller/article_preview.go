package controller

import (
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/router"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/yuin/goldmark"
)

func (c *Controller) GetArticlePreview(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return c.render(w, r, view.PageNotFound)
	}
	baseTemplates, err := c.Repository.SelectAllTemplate(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base template : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseTemplates {
		files[l.Name] = l.Content
	}
	layout, ok, err := c.Repository.SelectLayoutByID(r.Context(), a.LayoutID)
	if err != nil {
		return fmt.Errorf("cannot select layout %d : %w", a.LayoutID, err)
	}
	if !ok {
		return fmt.Errorf("cannot found layout %d : %w", a.LayoutID, err)
	}
	files[layout.Name] = layout.Content
	return renderPreview(w, files, a)
}

func renderPreview(w io.Writer, files map[string]string, pageData repository.Article) error {
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
		Content string
	}
	return tmpl.ExecuteTemplate(w, baseTemplate, PageData{
		Content: pageData.Content,
		Title:   pageData.Title,
	})
}
