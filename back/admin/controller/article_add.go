package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) PostArticleAdd(w http.ResponseWriter, r *http.Request) error {
	a, errors, err := form.ParseArticleAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return httpresponse.BadRequest(w, errors.Raw)
	}

	layouts, err := c.Repository.GetLayoutList(r.Context(), 1, 0)
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	if len(layouts) == 0 {
		return fmt.Errorf("need at leats one layout to create an article : %w", err)
	}

	slug, err := c.Repository.CreateArticle(r.Context(), repository.CreateArticleParams{
		Title:    a.Title,
		Author:   a.Author,
		Draft:    a.Draft.V,
		LayoutID: layouts[0].ID,
	})
	if err != nil {
		return fmt.Errorf("cannot create article : %w", err)
	}
	return view.ArticleCreated(w, view.ArticleAddData{
		Title:  a.Title,
		Author: a.Author,
		Draft:  a.Draft.V && a.Draft.Valid,
		Slug:   slug,
	})
}
