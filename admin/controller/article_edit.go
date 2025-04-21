package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticleEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:    a.Title,
		Author:   a.Author,
		Slug:     a.Slug,
		Content:  a.Content,
		Draft:    a.Draft,
		LayoutID: a.LayoutID,
		Layouts:  layoutSelector,
		Blocks:   blockSelector,
	}, view.ArticleEditError{}))
}

func (c *Controller) PostArticleEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	a, errors, err := form.ParseArticleEdit(r, c.Repository.CountLayoutByID)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	article.Title = a.Title
	article.Author = a.Author
	article.Content = a.Content
	article.Slug = a.Slug
	article.Draft = a.Draft
	article.LayoutID = a.LayoutID

	if !errors.HasError() {
		err := c.Repository.UpdateArticle(r.Context(), slug, article)
		if err != nil {
			return fmt.Errorf("cannot update %s article : %w", slug, err)
		}
	}
	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:    a.Title,
		Author:   a.Author,
		Slug:     a.Slug,
		Content:  a.Content,
		Draft:    a.Draft,
		LayoutID: a.LayoutID,
		Layouts:  layoutSelector,
		Blocks:   blockSelector,
	}, view.ArticleEditError(errors)))
}
