package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GeArticleEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:   a.Title,
		Author:  a.Author,
		Slug:    a.Slug,
		Content: a.Content,
		Draft:   a.Draft,
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

	a, errors, err := form.ParseArticleEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	article.Title = a.Title
	article.Author = a.Author
	article.Content = a.Content
	article.Slug = a.Slug
	article.Draft = a.Draft

	if !errors.HasError() {
		err := c.Repository.UpdateArticle(r.Context(), slug, article)
		if err != nil {
			return fmt.Errorf("cannot update %s article : %w", slug, err)
		}
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:   a.Title,
		Author:  a.Author,
		Slug:    a.Slug,
		Content: a.Content,
		Draft:   a.Draft,
	}, view.ArticleEditError(errors)))
}
