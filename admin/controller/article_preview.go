package controller

import (
	"app/pkg/router"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func (c *Controller) GetArticlePreview(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}
	baseLayouts, err := c.Repository.SelectBaseLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base layouts : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseLayouts {
		files[l.Name] = l.Content
	}
	pageLayoutName := "page/single.html"
	pageLayout, ok, err := c.Repository.SelectLayout(r.Context(), pageLayoutName)
	if err != nil {
		return fmt.Errorf("cannot select page layout %s : %w", pageLayoutName, err)
	}
	if !ok {
		return fmt.Errorf("cannot found page layout %s : %w", pageLayoutName, err)
	}
	files[pageLayoutName] = pageLayout.Content
	return renderPreview(w, files, a.Content)
}

func renderPreview(w io.Writer, files map[string]string, pageData string) error {
	const baseTemplate = "root"
	funcMap := template.FuncMap{}
	tmpl := template.New(baseTemplate).Funcs(funcMap)
	for name, content := range files {
		_, err := tmpl.New(name).Parse(content)
		if err != nil {
			return err
		}
	}
	type PageData struct {
		Content string
	}
	return tmpl.ExecuteTemplate(w, baseTemplate, PageData{Content: pageData})
}
