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
	baseLayouts, err := c.Repository.SelectAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base layouts : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseLayouts {
		files[l.Name] = l.Content
	}
	pageLayoutName := "page/single.html"
	pageLayout, ok, err := c.Repository.SelectPage(r.Context(), pageLayoutName)
	if err != nil {
		return fmt.Errorf("cannot select page layout %s : %w", pageLayoutName, err)
	}
	if !ok {
		return fmt.Errorf("cannot found page layout %s : %w", pageLayoutName, err)
	}
	files[pageLayoutName] = pageLayout.Content
	return renderPreview(w, files, a)
}

func renderPreview(w io.Writer, files map[string]string, pageData repository.Article) error {
	const baseTemplate = "baseof.html"
	if _, ok := files[baseTemplate]; !ok {
		return fmt.Errorf("base template %s not defined in _layout/", baseTemplate)
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
