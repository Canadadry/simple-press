package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticleAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.ArticleAdd(view.ArticleAddData{}, view.ArticleAddError{}))
}

func (c *Controller) PostArticleAdd(w http.ResponseWriter, r *http.Request) error {
	a, errors, err := form.ParseArticleAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return c.render(w, r, view.ArticleAdd(view.ArticleAddData(a), view.ArticleAddError(errors)))
	}

	pages, err := c.Repository.GetPageList(r.Context(), 1, 0)
	if err != nil {
		return fmt.Errorf("cannot select all pages : %w", err)
	}
	if len(pages) == 0 {
		return fmt.Errorf("need at leats one page to create an article : %w", err)
	}

	slug, err := c.Repository.CreateArticle(r.Context(), repository.CreateArticleParams{
		Title:    a.Title,
		Author:   a.Author,
		Draft:    a.Draft,
		LayoutID: pages[0].ID,
	})
	if err != nil {
		return fmt.Errorf("cannot create article : %w", err)
	}

	http.Redirect(w, r, "/admin/articles/"+slug+"/edit", http.StatusSeeOther)
	return nil
}
